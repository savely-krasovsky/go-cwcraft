package main

import (
	"errors"
)

func copyMap(m map[string]int) map[string]int {
	newMap := make(map[string]int)

	for key, value := range m {
		newMap[key] = value
	}

	return newMap
}

func findResourceByName(name string) (resource, error) {
	for _, res := range resources {
		if name == res.Name {
			return res, nil
		}
	}

	return resource{}, errors.New("not found")
}

func findResourceByID(ID string) (resource, error) {
	for _, res := range resources {
		if ID == res.ID {
			return res, nil
		}
	}

	return resource{}, errors.New("not found")
}

func getBasicsRecursively(basics quickMap, commands *[]command, recipe quickMap) (quickMap, []command) {
	for name, amount := range recipe {
		res, err := findResourceByName(name)
		if err != nil {
			// If can't find (for example unknown element, recipe or frag)
			basics[name] += amount
			continue
		}

		if res.Recipe == nil {
			// If it already basic
			basics[name] += amount
			continue
		} else {
			*commands = append(*commands, command{
				res.ID,
				res.Name,
				amount,
			})
		}

		// Copy (else we will change reference)
		recipe := copyMap(res.Recipe)

		// Multiple amount in recipe
		for recipeName, recipeAmount := range recipe {
			recipe[recipeName] = recipeAmount * amount
		}

		// Recursively go deeper
		getBasicsRecursively(basics, commands, recipe)
	}

	return basics, *commands
}