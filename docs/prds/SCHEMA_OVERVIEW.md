# Police Management System - Schema Overview

This document provides an overview of the implemented database schema for the police management system using Ent ORM.

## Entities

### 1. User

**Purpose**: Represents police officers and personnel in the system.

**Key Fields**:

- `id`: Unique identifier with prefix "user\_"
- `full_name`: Officer's full name
- `email`: Official email address
- `phone_number`: Contact number
- `rank`: Current rank/position
- `badge_number`: Official badge/service number
- `is_active`: Whether the user is currently active

**Relationships**:

- **Station**: Current station assignment (many-to-one)
- **Post**: Current post/rank (many-to-one)
- **Police Division**: Division the user is in charge of (one-to-one)
- **Hierarchical**: Superior/subordinate relationships with other users
- **Operations**: Operations performed on, by, or approved by this user
- **Permissions**: Permissions granted by or received by this user

### 2. Station

**Purpose**: Represents police stations and their basic information.

**Key Fields**:

- `id`: Unique identifier with prefix "station\_"
- `name`: Station name
- `station_code`: Unique station code
- `address`, `city`, `district`, `state`, `pincode`: Location details
- `jurisdiction_type`: Type of jurisdiction (rural, urban, suburban)
- `phone_number`, `email`: Contact information
- `last_audit_date`: Last audit date
- `feedback_score`: Citizen feedback score

**Relationships**:

- **Users**: Officers assigned to the station
- **SHO**: Station House Officer (one-to-one)
- **Operations**: Operations involving this station

### 3. StationDetail

**Purpose**: Extended details about station facilities and infrastructure.

**Key Fields**:

- `id`: Unique identifier with prefix "station*detail*"
- `station_type`: Type of station (Police Station, Police Outpost, etc.)
- `has_lockup`, `has_malkhana`, `has_parking`, `has_barracks`, `has_mess`: Facility flags
- `total_staff`, `sanctioned_strength`: Staffing information
- `infrastructure_status`: Overall infrastructure condition
- `special_facilities`, `equipment_details`: Additional details

**Relationships**:

- **Police Division**: Connected to a police division (one-to-one)

### 4. PoliceDivision

**Purpose**: Represents hierarchical police administrative divisions.

**Key Fields**:

- `id`: Unique identifier with prefix "police*division*"
- `name`: Division name
- `level`: Hierarchical level (state, zone, range, district, subdivision, circle, station)
- `code`: Official division code
- `address`: Division headquarters address
- `contact_number`, `email`: Contact information
- `jurisdiction_area`: Area under jurisdiction

**Relationships**:

- **Hierarchical**: Parent/child relationships with other divisions
- **Station Detail**: Additional details if this division is a station
- **Officer**: Officer in charge of this division
- **Operations**: Operations involving this division

### 5. Post

**Purpose**: Represents police ranks/positions with their hierarchical structure.

**Key Fields**:

- `id`: Unique identifier with prefix "post\_"
- `name`: Post/rank name (e.g., Constable, Inspector)
- `description`: Post responsibilities
- `hierarchy_level`: Numerical hierarchy level (higher = higher rank)
- `powers`: List of powers/authorities
- `max_subordinates`: Maximum direct subordinates
- `is_active`: Whether the post is currently active

**Relationships**:

- **Assignable By**: Posts that can assign officers to this post
- **Can Assign**: Posts that this post can assign officers to
- **Users**: Users currently holding this post
- **Operations**: Operations performed on this post

### 6. Operation

**Purpose**: Tracks all post operations (promotions, demotions, transfers, etc.).

**Key Fields**:

- `id`: Unique identifier with prefix "operation\_"
- `operation_type`: Type (promote, demote, transfer, remove, assign)
- `description`, `reason`: Operation details
- `before_details`, `after_details`: JSON fields for before/after state
- `effective_date`: When the operation becomes effective
- `is_approved`: Approval status
- `approved_at`: Approval timestamp

**Relationships**:

- **Target User**: User who is the target of the operation
- **Operated By**: User who performed the operation
- **Approved By**: User who approved the operation
- **Before/After Post**: Post before and after the operation
- **Before/After Station**: Station before and after the operation
- **Before/After Division**: Division before and after the operation

### 7. Permission

**Purpose**: Manages access control and permissions in the system.

**Key Fields**:

- `id`: Unique identifier with prefix "permission\_"
- `entity_type`, `entity_id`: Target entity for the permission
- `action`: Action type (read, write, execute, delete, assign, approve)
- `is_allowed`: Whether the action is allowed or denied
- `allow_entity_type`, `allow_entity_id`: Entity being granted/denied permission
- `conditions`: Additional conditions
- `valid_from`, `valid_until`: Permission validity period

**Relationships**:

- **Granted By**: User who granted the permission
- **Granted To User**: User who received the permission

## Key Features

### 1. Hierarchical Structure

- **Police Divisions**: Support parent/child relationships for administrative hierarchy
- **User Hierarchy**: Superior/subordinate relationships between officers
- **Post Hierarchy**: Numerical hierarchy levels for ranks

### 2. Audit Trail

- **Operations**: Complete tracking of all post operations with before/after states
- **Timestamps**: All entities have created_at/updated_at timestamps via BaseMixin
- **Approval Workflow**: Operations require approval with approval tracking

### 3. Access Control

- **Permissions**: Granular permission system for different actions
- **Entity-based**: Permissions can be granted on specific entities
- **Time-bound**: Permissions can have validity periods

### 4. Flexible Station Management

- **Basic Station Info**: Core station details in Station entity
- **Extended Details**: Additional facilities and infrastructure in StationDetail
- **SHO Assignment**: Dedicated Station House Officer assignment

### 5. UUID-based IDs

- All entities use UUID-based IDs with meaningful prefixes
- K-sortable UUIDs using ULID for better performance
- Consistent ID format across all entities

## Database Indexes

The schema includes optimized indexes for:

- **User**: Email, phone, badge number, rank, active status
- **Station**: Station code (unique)
- **Post**: Name (unique), hierarchy level, active status
- **Operation**: Operation type, approval status, effective date, created date
- **Permission**: Complex composite indexes for efficient permission lookups
- **PoliceDivision**: Hierarchical queries and level-based searches

## Usage Patterns

### 1. Officer Management

- Create users with posts and station assignments
- Track hierarchical relationships
- Manage promotions, transfers, and other operations

### 2. Station Administration

- Manage station details and facilities
- Assign SHOs and track station personnel
- Monitor station infrastructure and audits

### 3. Operations Tracking

- Record all personnel operations with full audit trail
- Implement approval workflows
- Track before/after states for all changes

### 4. Permission Management

- Grant granular permissions to users
- Control access to specific entities and actions
- Implement time-bound access controls

This schema provides a comprehensive foundation for a police management system with proper audit trails, hierarchical management, and flexible access control.
