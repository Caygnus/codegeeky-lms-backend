package internal

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/omkar273/codegeeky/internal/api/dto"
	"github.com/omkar273/codegeeky/internal/config"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/postgres"
	"github.com/omkar273/codegeeky/internal/repository"
	"github.com/omkar273/codegeeky/internal/service"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/shopspring/decimal"
)

// mockWebhookPublisher is a simple mock implementation for seeding
type mockWebhookPublisher struct{}

func (m *mockWebhookPublisher) PublishWebhook(ctx context.Context, event *types.WebhookEvent) error {
	// Do nothing for seeding - just return success
	return nil
}

func (m *mockWebhookPublisher) Close() error {
	// Do nothing for seeding
	return nil
}

// SeedData seeds the database with sample data using the service layer
func SeedData() error {
	// Load configuration
	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logger
	logger, err := logger.NewLogger(cfg)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	// Initialize database client
	entClient, err := postgres.NewEntClient(cfg, logger)
	if err != nil {
		return fmt.Errorf("failed to create database client: %w", err)
	}
	defer entClient.Close()

	// Create postgres client wrapper
	dbClient := postgres.NewClient(entClient, logger)

	// Create repository params
	repoParams := repository.RepositoryParams{
		Client: dbClient,
		Logger: logger,
		Config: cfg,
	}

	// Create repositories
	userRepo := repository.NewUserRepository(repoParams)
	discountRepo := repository.NewDiscountRepository(repoParams)
	paymentRepo := repository.NewPaymentRepository(repoParams)
	internshipRepo := repository.NewInternshipRepository(repoParams)
	internshipBatchRepo := repository.NewInternshipBatchRepository(repoParams)
	categoryRepo := repository.NewCategoryRepository(repoParams)
	internshipEnrollmentRepo := repository.NewInternshipEnrollmentRepository(repoParams)
	cartRepo := repository.NewCartRepository(repoParams)

	// Create service params
	serviceParams := service.ServiceParams{
		Logger:                   logger,
		Config:                   cfg,
		DB:                       dbClient,
		UserRepo:                 userRepo,
		DiscountRepo:             discountRepo,
		PaymentRepo:              paymentRepo,
		InternshipRepo:           internshipRepo,
		InternshipBatchRepo:      internshipBatchRepo,
		CategoryRepo:             categoryRepo,
		InternshipEnrollmentRepo: internshipEnrollmentRepo,
		CartRepo:                 cartRepo,
		WebhookPublisher:         &mockWebhookPublisher{}, // Use mock publisher for seeding
		HTTPClient:               nil,                     // Not needed for seeding
	}

	// Create services
	categoryService := service.NewCategoryService(serviceParams)
	internshipService := service.NewInternshipService(serviceParams)
	discountService := service.NewDiscountService(serviceParams)

	ctx := context.Background()

	log.Println("🌱 Starting database seeding...")

	// Step 1: Create Categories
	log.Println("📂 Creating categories...")
	categories, err := createCategories(ctx, categoryService)
	if err != nil {
		return fmt.Errorf("failed to create categories: %w", err)
	}
	log.Printf("✅ Created %d categories", len(categories))

	// Step 2: Create Internships with category relationships
	log.Println("💼 Creating internships...")
	internships, err := createInternships(ctx, internshipService, categories)
	if err != nil {
		return fmt.Errorf("failed to create internships: %w", err)
	}
	log.Printf("✅ Created %d internships", len(internships))

	// Step 3: Create Discounts
	log.Println("🎫 Creating discounts...")
	discounts, err := createDiscounts(ctx, discountService)
	if err != nil {
		return fmt.Errorf("failed to create discounts: %w", err)
	}
	log.Printf("✅ Created %d discounts", len(discounts))

	log.Println("🎉 Database seeding completed successfully!")
	log.Println("📊 Summary:")
	log.Printf("   - Categories: %d", len(categories))
	log.Printf("   - Internships: %d", len(internships))
	log.Printf("   - Discounts: %d", len(discounts))

	return nil
}

// createCategories creates sample categories using the category service
func createCategories(ctx context.Context, categoryService service.CategoryService) ([]*dto.CategoryResponse, error) {
	categories := []struct {
		name        string
		lookupKey   string
		description string
	}{
		{
			name:        "Web Development",
			lookupKey:   "web-development",
			description: "Full-stack web development internships covering frontend and backend technologies",
		},
		{
			name:        "Mobile Development",
			lookupKey:   "mobile-development",
			description: "Mobile app development for iOS and Android platforms",
		},
		{
			name:        "Data Science",
			lookupKey:   "data-science",
			description: "Data analysis, machine learning, and statistical modeling",
		},
		{
			name:        "DevOps",
			lookupKey:   "devops",
			description: "Infrastructure, deployment, and automation practices",
		},
		{
			name:        "UI/UX Design",
			lookupKey:   "ui-ux-design",
			description: "User interface and user experience design",
		},
		{
			name:        "Cybersecurity",
			lookupKey:   "cybersecurity",
			description: "Security testing, threat analysis, and secure coding practices",
		},
		{
			name:        "Cloud Computing",
			lookupKey:   "cloud-computing",
			description: "AWS, Azure, and Google Cloud platform development",
		},
		{
			name:        "Blockchain",
			lookupKey:   "blockchain",
			description: "Blockchain development and smart contract programming",
		},
	}

	var createdCategories []*dto.CategoryResponse

	for _, cat := range categories {
		req := &dto.CreateCategoryRequest{
			Name:        cat.name,
			LookupKey:   cat.lookupKey,
			Description: cat.description,
		}

		category, err := categoryService.Create(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to create category %s: %w", cat.name, err)
		}

		createdCategories = append(createdCategories, category)
	}

	return createdCategories, nil
}

// createInternships creates sample internships with category relationships using the internship service
func createInternships(ctx context.Context, internshipService service.InternshipService, categories []*dto.CategoryResponse) ([]*dto.InternshipResponse, error) {
	internships := []struct {
		title              string
		lookupKey          string
		description        string
		skills             []string
		level              types.InternshipLevel
		mode               types.InternshipMode
		durationInWeeks    int
		learningOutcomes   []string
		prerequisites      []string
		benefits           []string
		currency           string
		price              decimal.Decimal
		flatDiscount       *decimal.Decimal
		percentageDiscount *decimal.Decimal
		categoryIndex      int // Index in categories slice
	}{
		{
			title:           "Full-Stack Web Development",
			lookupKey:       "fullstack-web-dev",
			description:     "Learn modern full-stack web development with React, Node.js, and PostgreSQL. Build real-world applications and deploy them to production.",
			skills:          []string{"JavaScript", "React", "Node.js", "PostgreSQL", "Git", "Docker"},
			level:           types.InternshipLevelIntermediate,
			mode:            types.InternshipModeRemote,
			durationInWeeks: 12,
			learningOutcomes: []string{
				"Build responsive web applications",
				"Implement RESTful APIs",
				"Database design and optimization",
				"Deployment and CI/CD practices",
			},
			prerequisites: []string{"Basic JavaScript knowledge", "Understanding of HTML/CSS"},
			benefits: []string{
				"Portfolio of real projects",
				"Industry mentorship",
				"Job placement assistance",
				"Certificate of completion",
			},
			currency:      "USD",
			price:         decimal.NewFromInt(2999),
			categoryIndex: 0, // Web Development
		},
		{
			title:           "React Native Mobile Development",
			lookupKey:       "react-native-mobile",
			description:     "Master cross-platform mobile development with React Native. Build apps for both iOS and Android from a single codebase.",
			skills:          []string{"JavaScript", "React Native", "Redux", "Firebase", "App Store", "Google Play"},
			level:           types.InternshipLevelAdvanced,
			mode:            types.InternshipModeHybrid,
			durationInWeeks: 16,
			learningOutcomes: []string{
				"Cross-platform mobile development",
				"State management with Redux",
				"App store deployment",
				"Performance optimization",
			},
			prerequisites: []string{"React fundamentals", "JavaScript ES6+"},
			benefits: []string{
				"Published app on stores",
				"Performance optimization skills",
				"Real-world project experience",
			},
			currency:      "USD",
			price:         decimal.NewFromInt(3999),
			flatDiscount:  &[]decimal.Decimal{decimal.NewFromInt(500)}[0],
			categoryIndex: 1, // Mobile Development
		},
		{
			title:           "Machine Learning Fundamentals",
			lookupKey:       "ml-fundamentals",
			description:     "Introduction to machine learning with Python. Learn algorithms, data preprocessing, and model deployment.",
			skills:          []string{"Python", "Scikit-learn", "Pandas", "NumPy", "Matplotlib", "Jupyter"},
			level:           types.InternshipLevelBeginner,
			mode:            types.InternshipModeRemote,
			durationInWeeks: 10,
			learningOutcomes: []string{
				"Supervised and unsupervised learning",
				"Data preprocessing techniques",
				"Model evaluation and validation",
				"Feature engineering",
			},
			prerequisites: []string{"Basic Python programming", "High school mathematics"},
			benefits: []string{
				"ML project portfolio",
				"Kaggle competition experience",
				"Industry case studies",
			},
			currency:           "USD",
			price:              decimal.NewFromInt(2499),
			percentageDiscount: &[]decimal.Decimal{decimal.NewFromInt(15)}[0],
			categoryIndex:      2, // Data Science
		},
		{
			title:           "DevOps Engineering",
			lookupKey:       "devops-engineering",
			description:     "Learn modern DevOps practices including CI/CD, containerization, and cloud infrastructure management.",
			skills:          []string{"Docker", "Kubernetes", "AWS", "Jenkins", "Terraform", "Linux"},
			level:           types.InternshipLevelIntermediate,
			mode:            types.InternshipModeOnsite,
			durationInWeeks: 14,
			learningOutcomes: []string{
				"Container orchestration",
				"Infrastructure as Code",
				"CI/CD pipeline design",
				"Cloud architecture",
			},
			prerequisites: []string{"Basic Linux commands", "Understanding of networking"},
			benefits: []string{
				"Hands-on infrastructure experience",
				"AWS certification preparation",
				"Real deployment scenarios",
			},
			currency:      "USD",
			price:         decimal.NewFromInt(3499),
			categoryIndex: 3, // DevOps
		},
		{
			title:           "UI/UX Design Masterclass",
			lookupKey:       "ui-ux-masterclass",
			description:     "Master the principles of user interface and user experience design. Create beautiful, functional, and user-friendly designs.",
			skills:          []string{"Figma", "Adobe XD", "Sketch", "Prototyping", "User Research", "Design Systems"},
			level:           types.InternshipLevelAdvanced,
			mode:            types.InternshipModeHybrid,
			durationInWeeks: 12,
			learningOutcomes: []string{
				"User research and personas",
				"Wireframing and prototyping",
				"Design system creation",
				"Usability testing",
			},
			prerequisites: []string{"Basic design principles", "Familiarity with design tools"},
			benefits: []string{
				"Professional design portfolio",
				"Industry design challenges",
				"Design system certification",
			},
			currency:      "USD",
			price:         decimal.NewFromInt(2799),
			categoryIndex: 4, // UI/UX Design
		},
		{
			title:           "Cybersecurity Fundamentals",
			lookupKey:       "cybersecurity-fundamentals",
			description:     "Learn essential cybersecurity concepts including ethical hacking, secure coding, and threat analysis.",
			skills:          []string{"Python", "Wireshark", "Metasploit", "OWASP", "Network Security", "Cryptography"},
			level:           types.InternshipLevelBeginner,
			mode:            types.InternshipModeRemote,
			durationInWeeks: 8,
			learningOutcomes: []string{
				"Network security analysis",
				"Vulnerability assessment",
				"Secure coding practices",
				"Incident response",
			},
			prerequisites: []string{"Basic networking knowledge", "Programming fundamentals"},
			benefits: []string{
				"Security certification prep",
				"Real-world security challenges",
				"Industry security tools",
			},
			currency:      "USD",
			price:         decimal.NewFromInt(1999),
			categoryIndex: 5, // Cybersecurity
		},
	}

	var createdInternships []*dto.InternshipResponse

	for _, intern := range internships {
		// Get category ID
		if intern.categoryIndex >= len(categories) {
			return nil, fmt.Errorf("invalid category index %d for internship %s", intern.categoryIndex, intern.title)
		}
		categoryID := categories[intern.categoryIndex].ID

		req := &dto.CreateInternshipRequest{
			Title:              intern.title,
			LookupKey:          intern.lookupKey,
			Description:        intern.description,
			Skills:             intern.skills,
			Level:              intern.level,
			Mode:               intern.mode,
			DurationInWeeks:    intern.durationInWeeks,
			LearningOutcomes:   intern.learningOutcomes,
			Prerequisites:      intern.prerequisites,
			Benefits:           intern.benefits,
			Currency:           intern.currency,
			Price:              intern.price,
			FlatDiscount:       intern.flatDiscount,
			PercentageDiscount: intern.percentageDiscount,
			CategoryIDs:        []string{categoryID},
		}

		internship, err := internshipService.Create(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to create internship %s: %w", intern.title, err)
		}

		createdInternships = append(createdInternships, internship)
	}

	return createdInternships, nil
}

// createDiscounts creates sample discounts using the discount service
func createDiscounts(ctx context.Context, discountService service.DiscountService) ([]*dto.DiscountResponse, error) {
	discounts := []struct {
		code          string
		description   string
		discountType  types.DiscountType
		discountValue decimal.Decimal
		validFrom     time.Time
		validUntil    *time.Time
		isActive      bool
		maxUses       *int
		minOrderValue *decimal.Decimal
		isCombinable  bool
	}{
		{
			code:          "WELCOME20",
			description:   "Welcome discount for new students",
			discountType:  types.DiscountTypePercentage,
			discountValue: decimal.NewFromInt(20),
			validFrom:     time.Now(),
			validUntil:    &[]time.Time{time.Now().AddDate(0, 3, 0)}[0], // 3 months
			isActive:      true,
			maxUses:       &[]int{100}[0],
			minOrderValue: &[]decimal.Decimal{decimal.NewFromInt(1000)}[0],
			isCombinable:  false,
		},
		{
			code:          "FLAT500",
			description:   "Flat discount for premium courses",
			discountType:  types.DiscountTypeFlat,
			discountValue: decimal.NewFromInt(500),
			validFrom:     time.Now(),
			validUntil:    &[]time.Time{time.Now().AddDate(0, 6, 0)}[0], // 6 months
			isActive:      true,
			maxUses:       &[]int{50}[0],
			minOrderValue: &[]decimal.Decimal{decimal.NewFromInt(2000)}[0],
			isCombinable:  true,
		},
		{
			code:          "SUMMER15",
			description:   "Summer special discount",
			discountType:  types.DiscountTypePercentage,
			discountValue: decimal.NewFromInt(15),
			validFrom:     time.Now(),
			validUntil:    &[]time.Time{time.Now().AddDate(0, 2, 0)}[0], // 2 months
			isActive:      true,
			maxUses:       &[]int{200}[0],
			minOrderValue: &[]decimal.Decimal{decimal.NewFromInt(500)}[0],
			isCombinable:  false,
		},
		{
			code:          "BULK25",
			description:   "Bulk purchase discount",
			discountType:  types.DiscountTypePercentage,
			discountValue: decimal.NewFromInt(25),
			validFrom:     time.Now(),
			validUntil:    &[]time.Time{time.Now().AddDate(1, 0, 0)}[0], // 1 year
			isActive:      true,
			maxUses:       &[]int{25}[0],
			minOrderValue: &[]decimal.Decimal{decimal.NewFromInt(5000)}[0],
			isCombinable:  false,
		},
		{
			code:          "EARLYBIRD100",
			description:   "Early bird flat discount",
			discountType:  types.DiscountTypeFlat,
			discountValue: decimal.NewFromInt(100),
			validFrom:     time.Now(),
			validUntil:    &[]time.Time{time.Now().AddDate(0, 1, 0)}[0], // 1 month
			isActive:      true,
			maxUses:       &[]int{75}[0],
			minOrderValue: &[]decimal.Decimal{decimal.NewFromInt(1500)}[0],
			isCombinable:  true,
		},
	}

	var createdDiscounts []*dto.DiscountResponse

	for _, disc := range discounts {
		req := &dto.CreateDiscountRequest{
			Code:          disc.code,
			Description:   disc.description,
			DiscountType:  disc.discountType,
			DiscountValue: disc.discountValue,
			ValidFrom:     &disc.validFrom,
			ValidUntil:    disc.validUntil,
			IsActive:      &disc.isActive,
			MaxUses:       disc.maxUses,
			MinOrderValue: disc.minOrderValue,
			IsCombinable:  disc.isCombinable,
			Metadata: types.Metadata{
				"created_by": "seed_script",
				"purpose":    "sample_data",
			},
		}

		discount, err := discountService.Create(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to create discount %s: %w", disc.code, err)
		}

		createdDiscounts = append(createdDiscounts, discount)
	}

	return createdDiscounts, nil
}
