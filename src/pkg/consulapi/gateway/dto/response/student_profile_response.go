package response

type StudentProfileResponse struct {
	StudentInformation StudentInformation `json:"student_information"`
}

type StudentInformation struct {
	DOB               string `json:"dob"`
	Gender            uint   `json:"gender"`
	StudyLevel        uint   `json:"study_level"`
	MinWaterMustDrink uint   `json:"min_water_must_drink"`
	Description       string `json:"description"`
	Mode              string `json:"mode"`
}
