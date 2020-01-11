package main

import (
	"fmt"
	"sync"
	"time"
)

//global 

var owner []int //array of who is using mutex
var claim []int //array of who has a claim on a mutex
var debug = 0 // variable for printing debug info
var delay = 1 // variable for ensuring deadlock
var claimCheck = 0 // seperate debug variable for solution 3
var claimMutex = sync.Mutex{} //solution 4 mutex
var arbiter = sync.Mutex{} // solution 2 mutex


type utensil struct {
	mutex sync.Mutex
	num int
}


func hungryPhilo(philNum int, right,left *utensil,soln int){ // function that picks solution to problem

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


func dibs(right, left, philNum int) int{ // this is my version of message passing through philosophers
// should probably be more of the pipe/channel system for go routines but I already had the infrastructure

  claimMutex.Lock()

	if owner[left] == -1{  //if utensil clean

		 if claim[left] == philNum{ // and we have a claim on the left mutex

		 	if claim[right] == philNum { // and we have a claim on the right mutex
				owner[left] = philNum
				owner[right] = philNum
				claimMutex.Unlock()
				return 1 // we claim the resource

			}else if claim[right]!=-1{ //and someone has dibs on the right

				claimMutex.Unlock()
				return 0 //leave resource alone
			} else { //and there is no claim
				claim[right] = philNum
				claimMutex.Unlock()
				return 0 // then claim
			}
		 } else if claim[left]!=-1 {// and someone else does
			 claimMutex.Unlock()
			 return 0

		 } else {
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
	if claimCheck == 1 {
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


func debugAndDelayPrint(text string,a ...int){ // this prints statements if debug is on and adds delay if variable is set. 
	if debug == 1 {
		fmt.Printf(text,a[0],a[1])
	}
	if delay == 1 {
		time.Sleep(1*time.Millisecond)
	}
}


func eat(philNum, soln int,right,left *utensil){ // this is the main utensil taking and releasing code


	left.mutex.Lock()
	owner[left.num] = philNum
	debugAndDelayPrint("PH %v: acquired mutex %v...\n",philNum,left.num)

	right.mutex.Lock()
	owner[right.num] = philNum
	debugAndDelayPrint("PH %v: acquired mutex %v...\n" ,philNum, right.num)

	if soln == 2{
		arbiter.Unlock()
		debugAndDelayPrint("PH %v: obtained %v utensils, releasing arbitor\n", philNum, 2 )
	}

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


func initOwner(){ // initializes mutex owner slice
	for i := 0; i< len(owner); i++{
		owner[i] = -1
	}
}


func initClaim(){ // initializes claim to mutex slice
	for i := 0; i< len(claim); i++{
		claim[i] = -1
	}
}


func main() {

	var timer = 0
	n := 3
  	solnNumber := 0 // 0 is no fix, 1 is left hand soln, 2 is arbiter soln, 3 is chandy-ish soln
	owner = make([]int,n )
	claim = make([]int,n )
	initOwner()
	initClaim()
	utensils := make([]*utensil, n)
	for i :=0; i< n; i++ { // init utensils
		utensils[i] = new(utensil)
		utensils[i].num = i
	}

	for i:= 0; i<n; i++ { // init philosophers
		go hungryPhilo(i,utensils[i],utensils[(i+1)%n],solnNumber)
	}

	if debug == 1{
		fmt.Printf("Debug on\n")
	}
	if delay == 1{
		fmt.Printf("Delay on for ensured deadlock\n")
	}
	for{
		timer++
		time.Sleep(15* time.Second)
		fmt.Printf("MAIN: waited for %v seconds\n", 15*timer)
		
		if debug == 1 {
			fmt.Printf("Owner of mutexes: %v\n", owner)
		}
	}
}
