package model

type Method interface {
	// GetMaxVersionCount
	//
	// select COUNT(*) from users group by version order by version desc limit 1
	GetMaxVersionCount() (int64, error)
}
