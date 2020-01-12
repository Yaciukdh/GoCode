package main

import (
	"fmt"
	"sync"
	"time"
)

var owner []int
var claim []int
var foodEaten []int
var debug = 0
var delay = 0
var claimCheck = 0
var foodPrint = 1
var claimMutex = sync.Mutex{}
var arbiter = sync.Mutex{}

type utensil struct {
	mutex sync.Mutex
	num int
}

func hungryPhilo(philNum int, right,left *utensil,soln int){// function that picks solution to problem

	if soln == 0 {
		for { // original problem, no solution implemented
			eat(philNum,soln, right, left)
		}
	} else if soln == 1 { //left hand solution
		for {
			if philNum != 0 {
				eat(philNum,soln, right, left)
			} else {
				eat(philNum,soln, left, right)
			}
		}
	}else if soln == 2 { //permission from server/arbiter solution
		for{
			permissionToEat := ask(right.num,left.num)
			if permissionToEat == 1 {
				eat(philNum,soln, right, left)
			}

		}
	}else if soln == 3{ // augmented chandy misra solution
		for{
			cv := dibs(right.num,left.num, philNum)
			if cv == 1 {
				eat(philNum,soln, right, left)
				checkClaims()
			}
		}
	}

}


func dibs(right, left, philNum int) int{// this is my version of message passing through philosophers
// should probably be more of the pipe/channel system for go routines but I already had this infrastructure

	claimMutex.Lock()

	if owner[left] == -1{  //if utensil clean

		 if claim[left] == philNum{ // and we have a claim on the left mutex

		 	if claim[right] == philNum { // and we have a claim on the right mutex
				owner[left] = philNum
				owner[right] = philNum
				claimMutex.Unlock()
				return 1 // we claim the resource

			}else if claim[right]!=-1{ //and someone has dibs on the right
				claim[left] = -1 //give up claim to left
				claimMutex.Unlock()
				return 0
			} else { //and there is no claim
				claim[right] = philNum
				claimMutex.Unlock()
				return 0 // then claim
			}
		 } else if claim[left]!=-1 {// and someone else does
			 claimMutex.Unlock()
			 return 0

		 } else { // no one has claim on mutex
		 	claim[left] = philNum
		 	claimMutex.Unlock()
			return 0
		 }
	}else{ //utensil is in use

		claimMutex.Unlock()
		return 0
	}

}

func checkClaims(){ // this just prints claims if the variavle claimsCheck is set

	if claimCheck ==1 {
		claimMutex.Lock()
		fmt.Println(claim)
		claimMutex.Unlock()
	}
}


func removeDibs(right,left int) { // this removes the owner's claim to the utensils

	claimMutex.Lock()
	claim[right]= -1
	claim[left] = -1
	claimMutex.Unlock()
}

func ask(right, left int) int{ // this is for soln 2, where you ask a server to distribute utensils

	arbiter.Lock()
	if owner[right] == -1 && owner[left] == -1 {
		debugAndDelayPrint("asking for %v and %v\n", right, left)
		return 1
	} else{
		arbiter.Unlock()
	}
	return 0
}

func debugAndDelayPrint(text string,a ...int){//this prints statements if debug is on and adds delay if variable is set

	if debug == 1 {
		fmt.Printf(text,a[0],a[1])
	}
	if delay == 1 {
		time.Sleep(1*time.Millisecond)
	}
}

func eat(philNum, soln int,right,left *utensil){// this is the main utensil taking and releasing code

	left.mutex.Lock()
	owner[left.num] = philNum
	debugAndDelayPrint("PH %v: acquired mutex %v...\n",philNum,left.num)

	right.mutex.Lock()
	owner[right.num] = philNum
	debugAndDelayPrint("PH %v: acquired mutex %v...\n" ,philNum, right.num)

	if soln == 2 {
		arbiter.Unlock()
		debugAndDelayPrint("PH %v: obtained %v utensils, releasing arbitor\n", philNum, 2 )
	}

	foodEaten[philNum]++
	fmt.Printf("PH %v: eating for 1 seconds...\n",philNum)
	time.Sleep(2 * time.Second)

	right.mutex.Unlock()
	owner[right.num] = -1
	debugAndDelayPrint("PH %v: releasing mutex %v...\n",philNum,left.num)

	left.mutex.Unlock()
	owner[left.num] = -1
	debugAndDelayPrint("PH %v: releasing mutex %v...\n",philNum,right.num)

	if soln == 3{
		removeDibs(right.num,left.num)
	}

	fmt.Printf("PH %v: is full, thinking for 1 second...\n",philNum)
	time.Sleep(1 * time.Second)

}

func initVars(){ // initializes variables

	for i := 0; i< len(owner); i++{
		owner[i] = -1
		claim[i] = -1
	}

}


func main() {

	var timer = 0
	n := 8

	owner = make( []int,n )
	claim = make( []int,n )
	foodEaten = make( []int,n )
	utensils := make( []*utensil,n )

	initVars()

	for i :=0; i< n; i++ {
		utensils[i] = new(utensil)
		utensils[i].num = i
	}

	for i:= 0; i<n; i++ {
		go hungryPhilo(i,utensils[i],utensils[(i+1)%n],1)
	}

	if debug == 1{
		fmt.Printf("Debug on\n")
	}
	if delay == 1{
		fmt.Printf("Delay on for ensured deadlock\n")
	}
	for{
		timer++
		time.Sleep(15*time.Second)
		fmt.Printf("MAIN: waited for %v seconds\n", 15*timer)

		if debug == 1 {
			fmt.Printf("Owner of mutexes: %v\n", owner)
		}
		if foodPrint == 1{
			fmt.Printf("Food eaten: %v\n", foodEaten)
		}
	}
}
