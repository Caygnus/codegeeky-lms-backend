package internship

import (
	"time"

	"github.com/omkar273/codegeeky/ent"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
)

// InternshipBatch is the model entity for the InternshipBatch schema.
type InternshipBatch struct {
	ID             string                      `json:"id,omitempty"`
	InternshipID   string                      `json:"internship_id,omitempty"`
	Name           string                      `json:"name,omitempty"`
	Description    string                      `json:"description,omitempty"`
	StartDate      time.Time                   `json:"start_date,omitempty"`
	EndDate        time.Time                   `json:"end_date,omitempty"`
	BatchStatus    types.InternshipBatchStatus `json:"batch_status,omitempty"`
	types.Metadata `json:"metadata,omitempty"`
	types.BaseModel
}

func (b *InternshipBatch) FromEnt(ent *ent.InternshipBatch) *InternshipBatch {
	return &InternshipBatch{
		ID:           ent.ID,
		InternshipID: ent.InternshipID,
		Name:         ent.Name,
		Description:  ent.Description,
		StartDate:    ent.StartDate,
		EndDate:      ent.EndDate,
		BatchStatus:  types.InternshipBatchStatus(ent.BatchStatus),
		Metadata:     types.MetadataFromEnt(ent.Metadata),
		BaseModel: types.BaseModel{
			Status:    types.Status(ent.Status),
			CreatedAt: ent.CreatedAt,
			UpdatedAt: ent.UpdatedAt,
			CreatedBy: ent.CreatedBy,
			UpdatedBy: ent.UpdatedBy,
		},
	}
}

func (b *InternshipBatch) FromEntList(ents []*ent.InternshipBatch) []*InternshipBatch {
	return lo.Map(ents, func(ent *ent.InternshipBatch, _ int) *InternshipBatch {
		return b.FromEnt(ent)
	})
}
