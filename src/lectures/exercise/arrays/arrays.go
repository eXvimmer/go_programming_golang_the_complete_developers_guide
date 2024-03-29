//--Summary:
//  Create a program that can store a shopping list and print out information
//  about the list.
//
//--Requirements:
//* Using an array, create a shopping list with enough room
//  for 4 products
//  - Products must include the price and the name
//* Insert 3 products into the array
//* Print to the terminal:
//  - The last item on the list
//  - The total number of items
//  - The total cost of the items
//* Add a fourth product to the list and print out the
//  information again

package main

import "fmt"

type Product struct {
	Name  string
	Price float64
}

func printStats(list [4]Product) {
	var cost float64
	totalItems := 0
	for i := 0; i < len(list); i++ {
		item := list[i]
		cost += item.Price

		if item.Name != "" {
			totalItems++
		}
	}

	fmt.Println("last item on the list:", list[totalItems-1])
	fmt.Println("Total items:", totalItems)
	fmt.Println("Total cost:", cost)
}

func main() {
	var sl [4]Product
	sl[0] = Product{Name: "milk", Price: 3.99}
	sl[1] = Product{Name: "egg", Price: 9.99}
	sl[2] = Product{Name: "cheese", Price: 4.99}

	printStats(sl)
}
