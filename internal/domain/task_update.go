package domain

// TaskUpdate represents the fields that can be updated for a task
type TaskUpdate struct {
	Title       *string `bson:"title,omitempty"`
	Description *string `bson:"description,omitempty"`
	Status      *string `bson:"status,omitempty"`
}
