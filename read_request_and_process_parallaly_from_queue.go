package api

import (
	"fmt"
	"net/http"
	"strconv"
)

var queue = make(chan int, 100)

func RunQueue() {
	
	go func() {
		
		//go routine to print queue data
		for data := range queue {
			fmt.Printf("1st approach %d \n", data)
		}

		// select {
		// case data := <- queue:
		// 	fmt.Println(data)

		// default:

		// 	fmt.Println("default case")
		// }

	
	}()
	http.HandleFunc("/queue", processQueue)

	http.ListenAndServe(":8080", nil)

	
}

func processQueue(w http.ResponseWriter, r *http.Request) {

	// Get query parameters from the URL
	queryParams := r.URL.Query()

	// Access specific parameters
	data := queryParams.Get("element")
	element, err := strconv.Atoi(data)


	if err != nil {
		fmt.Fprintln(w, "Badrequest. Element is not a number")
		return 
	}
	queue <- element

	fmt.Fprintln(w, "Reading from qeueue")

}

//Custom Queue
type CustomQueue struct {
	elements []int
}
func(q *CustomQueue) Enqueue(element int) {
	q.elements = append(q.elements, element)
}

func(q *CustomQueue) Dequeue() int{
	if len(q.elements) == 0 {
		return -1 //no data in qeueue
	}
	
	firstElement := q.elements[0]
	//Remove first element and return it
	q.elements = q.elements[1:] //takes second element onwards
	return firstElement
}
