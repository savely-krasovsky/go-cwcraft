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

func findResource(name string) (resource, error) {
	for _, resource := range resources {
		if name == resource.Name {
			return resource, nil
		}
	}

	return resource{}, errors.New("not found")
}

func getBasicRecursively(basic map[string]int, recipe map[string]int) map[string]int {
	for name, amount := range recipe {
		resource, err := findResource(name)
		if err != nil {
			// If can't find (for example unknown element, recipe or frag)
			basic[name] += amount
			continue
		}

		if resource.Recipe == nil {
			// If it already basic
			basic[name] += amount
			continue
		}

		// Copy (else we will change reference)
		recipe := copyMap(resource.Recipe)

		// Multiple amount in recipe
		for recipeName, recipeAmount := range recipe {
			recipe[recipeName] = recipeAmount * amount
		}

		// Recursively go deeper
		getBasicRecursively(basic, recipe)
	}

	return basic
}