package entity

type SPdf struct {
	ID        uint64 `gorm:"primary_key;auto_increment;"`
	PdfName   string `gorm:"column:pdf_name;not null;"`
	Folder    string `gorm:"column:folder;not null;"`
	Key       string `gorm:"column:key;not null;unique;"`
	Extension string `gorm:"column:extension;not null;"`
}

func (SPdf) TableName() string {
	return "s_pdf"
}
