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
	for _, res := range resources {
		if name == res.Name {
			return res, nil
		}
	}

	return resource{}, errors.New("not found")
}

func getBasicsRecursively(basics quickMap, commands *[]command, recipe quickMap) (quickMap, []command) {
	for name, amount := range recipe {
		res, err := findResource(name)
		if err != nil {
			// if can't find (for example unknown element, recipe or frag)
			basics[name] += amount
			continue
		}

		if res.Recipe == nil {
			// if it already basic
			basics[name] += amount
			continue
		} else {
			// let's say that -1 means "there is commands with this resource"
			comIndex := -1

			// find index of resource
			for i, c := range *commands {
				if name == c.Name {
					comIndex = i
				}
			}

			// have found? add it to already existing command
			if comIndex != -1 {
				(*commands)[comIndex].Amount += amount
				(*commands)[comIndex].CommandManaCost += res.ManaCost * amount
			} else {
				// else just add new
				*commands = append(*commands, command{
					res.ID,
					res.Name,
					amount,
					res.ManaCost * amount,
				})
			}
		}

		// copy (else we will change reference)
		recipe := copyMap(res.Recipe)

		// multiple amount in recipe
		for recipeName, recipeAmount := range recipe {
			recipe[recipeName] = recipeAmount * amount
		}

		// recursively go deeper
		getBasicsRecursively(basics, commands, recipe)
	}

	return basics, *commands
}
