package skills

type Skill struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Source      string `json:"source"`
	Version     string `json:"version"`
}

type Service struct {
	skills map[string]Skill
}

func NewService() *Service {
	return &Service{
		skills: make(map[string]Skill),
	}
}

func (s *Service) ListSkills() ([]Skill, error) {
	var result []Skill
	for _, skill := range s.skills {
		result = append(result, skill)
	}
	return result, nil
}

func (s *Service) GetSkill(name string) (*Skill, error) {
	skill, ok := s.skills[name]
	if !ok {
		return nil, nil
	}
	return &skill, nil
}
