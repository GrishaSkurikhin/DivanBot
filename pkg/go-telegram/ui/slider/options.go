package slider

type Option func(s *Slider)

func Button(buttonText string, buttonData string) Option {
	return func(s *Slider) {
		s.buttonsText = append(s.buttonsText, buttonText)
		s.buttonsData = append(s.buttonsData, buttonData)
	}
}
