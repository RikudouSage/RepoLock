package helper

import "regexp"

func RegexNamedSubmatch(regex *regexp.Regexp, target string) map[string]string {
	match := regex.FindStringSubmatch(target)
	if match == nil {
		return nil
	}

	groups := make(map[string]string, len(match))

	for i, name := range regex.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}

		groups[name] = match[i]
	}

	return groups
}
