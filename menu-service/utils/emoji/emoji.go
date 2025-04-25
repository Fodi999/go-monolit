package emoji

import "strings"

func GenerateEmoji(name string) string {
	n := strings.ToLower(name)

	switch {
	case strings.Contains(n, "огурец"):
		return "🥒"
	case strings.Contains(n, "помидор"):
		return "🍅"
	case strings.Contains(n, "рис"):
		return "🍚"
	case strings.Contains(n, "суши"):
		return "🍣"
	case strings.Contains(n, "лосось"), strings.Contains(n, "семга"):
		return "🐟"
	case strings.Contains(n, "соус"):
		return "🥣"
	case strings.Contains(n, "сыр"):
		return "🧀"
	case strings.Contains(n, "авокадо"):
		return "🥑"
	case strings.Contains(n, "морковь"):
		return "🥕"
	default:
		return "🍽️"
	}
}
