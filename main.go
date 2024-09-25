package main

import (
	"bufio"
	"context"
	// "encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"

	// "encoding/json"
	"log"

	"github.com/joho/godotenv"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
type(
todoModel struct{
	ID primitive.ObjectID `bson:"_id,omitEmpty"`
	Title string `bson:"title"`
	Completed bool `bson:"completed"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`

}
todo struct{
	ID string `json:"id"`
	Title string `json:"title"`
	Completed bool `json:"completed"`
	CreatedAt time.Time `json:"createdAt`
	UpdatedAt time.Time `json:"updatedAt"`
}
)

type Task struct {
	TaskName string
	completed bool
}
var task[] Task
func addTask(task string, coll *mongo.Collection) {
	newTask := Task{TaskName: task, completed: false}
    fmt.Println(task)
    tm :=todoModel{
		ID : primitive.NewObjectID(),
		Title: newTask.TaskName,
		Completed: false,
		CreatedAt : time.Now(),
	}
    fmt.Println(tm.ID)
	result, err := coll.InsertOne(context.TODO(), tm)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Task Added",result)
	
}

func listTasks(coll *mongo.Collection){
	// for i, task := range tasks {
	// 	status := "Pending"
	// 	if task.completed {
	// 		status = "Done"
	// 	}
	// 	fmt.Printf("%d. %s [%s]\n", i+1, task.TaskName, status)
	// }
//    var todo []todoModel
   todos := []todoModel{}
   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()

   // Fetch all documents
   cursor, err := coll.Find(ctx, bson.D{})
   if err != nil {
	  log.Println(err)
   }
   defer cursor.Close(ctx)
   if err := cursor.All(ctx, &todos); err != nil {
	log.Println(err)
}
todoList := []todo{}
for _,t := range todos{
	todoList = append(todoList, todo{
		ID: t.ID.String(),
		Title: t.Title,
		Completed: t.Completed,
		CreatedAt: t.CreatedAt,
	})
}
for _,t :=range todoList{
	fmt.Println(t)
}
}

func markCompleted(index string, coll *mongo.Collection){
	// if index >= 1 && index <= len(tasks){
	// 	tasks[index-1].completed = true
	// } else {
	// 	fmt.Println("Invalid Index")
	// }
	 // Convert the string ObjectID to a MongoDB ObjectID
	 objectID, err := primitive.ObjectIDFromHex(index)
	 if err != nil {
		 log.Fatalf("Invalid ObjectID: %v\n", err)
	 }
 
	 // Specify the filter to find the document by ObjectID
	 filter := bson.M{"_id": objectID}

     istLocation, err := time.LoadLocation("Asia/Kolkata")
    if err != nil {
        fmt.Println("Error loading location:", err)
        return
    }
	 // Define the update using the $set operator to modify fields
	 update := bson.M{
		 "$set": bson.M{
			//  "title":     "Updated Title",      // Replace this with the new title
			 "completed": true,                 // Example of updating a boolean field
			 "updatedAt": time.Now().In(istLocation),           // Set the updated time
		 },
	 }
 
	 // Update the document
	 updateResult, err := coll.UpdateOne(context.TODO(), filter, update)
	 if err != nil {
		 log.Fatal(err)
	 }
 
	 // Print the result of the update operation
	 if updateResult.ModifiedCount > 0 {
		 fmt.Printf("Updated %v document(s)\n", updateResult.ModifiedCount)
	 } else {
		 fmt.Println("No documents were updated")
	 }
}

func editTask(index string, newName string, coll *mongo.Collection){
	// if index >= 1 && index <= len(tasks){
	// 	tasks[index-1].TaskName = newName
	// 	fmt.Println("Task edited sucessfully")
		
	// } else {
	// 	fmt.Println("Invalid Index")
	// }
	objectID, err := primitive.ObjectIDFromHex(index)
	 if err != nil {
		 log.Fatalf("Invalid ObjectID: %v\n", err)
	 }
 
	 // Specify the filter to find the document by ObjectID
	 filter := bson.M{"_id": objectID}

     istLocation, err := time.LoadLocation("Asia/Kolkata")
    if err != nil {
        fmt.Println("Error loading location:", err)
        return
    }
	 // Define the update using the $set operator to modify fields
	 update := bson.M{
		 "$set": bson.M{
			 "title":     newName,      // Replace this with the new title
			//  "completed": true,                 // Example of updating a boolean field
			 "updatedAt": time.Now().In(istLocation),           // Set the updated time
		 },
	 }
 
	 // Update the document
	 updateResult, err := coll.UpdateOne(context.TODO(), filter, update)
	 if err != nil {
		 log.Fatal(err)
	 }
 
	 // Print the result of the update operation
	 if updateResult.ModifiedCount > 0 {
		 fmt.Printf("Updated %v document(s)\n", updateResult.ModifiedCount)
	 } else {
		 fmt.Println("No documents were updated")
	 }
}

func deleteTask(index string, coll *mongo.Collection){
	// if index >= 1 && index <= len(tasks){
	// 	tasks = append(tasks[:index-1], tasks[index:]...)
	// 	fmt.Println("Task deleted sucessfully")
		
	// } else {
	// 	fmt.Println("Invalid Index")
	// }

	 // Create a context with a timeout
	 ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	 defer cancel()
 
	 // Filter by ObjectID to delete a specific document
	 id, err := primitive.ObjectIDFromHex(index) // Replace with the actual ObjectID
	 filter := bson.M{"_id": id}
     if err!=nil{
		log.Println("id not exist", err)
	 }
	 // Delete the document
	 deleteResult, err := coll.DeleteOne(ctx, filter)
	 if err != nil {
		 log.Fatal(err)
	 }
 
	 // Output the result
	 fmt.Printf("Deleted %v document(s)\n", deleteResult.DeletedCount)
}
func main(){
  err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }
  uri := os.Getenv("MONGODB_URI")
  fmt.Println(uri)
  
  clientOptions := options.Client().ApplyURI(uri)
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }

    // Ensure disconnection on function exit
    defer func() {
        if err = client.Disconnect(context.TODO()); err != nil {
            log.Fatalf("Failed to disconnect from MongoDB: %v", err)
        }
    }()

    // Check the connection
    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatalf("Ping to MongoDB failed: %v", err)
    }

    fmt.Println("Connected to MongoDB!")

  coll := client.Database("todo_list").Collection("vp")

//   var indexInput int
  var taskInput, newTaskInput string
  choices := []string{"Add Task", "List Tasks", "Mark as Complete", "Edit Task", "Delete Task", "Exit"}
  
  fmt.Println("Options")
  for i, option := range choices{
	fmt.Printf("%d. %s\n",i+1, option)
  }
  
  scanner := bufio.NewScanner(os.Stdin)

  for{
	fmt.Print("Enter choice (1,2,3,4,5,6): ")
	scanner.Scan()
	input := scanner.Text()

	choice, err := strconv.Atoi(input)
	
	if err != nil {
		fmt.Println("Invalid Choice")
		continue
	}
	switch choice {
	case 1:
	  fmt.Print("Enter task: ")
	  scanner.Scan()
	  taskInput = scanner.Text()
	  addTask(taskInput,coll)
	case 2:
		listTasks(coll)
	case 3:
		fmt.Print("Enter index: ")
		scanner.Scan()

		id := scanner.Text()
		markCompleted(id, coll)
	case 4:
		fmt.Print("Enter index: ")
		scanner.Scan()
		id:= scanner.Text()
		fmt.Print("Enter task: ")
	    scanner.Scan()
	    newTaskInput = scanner.Text()
	    editTask(id, newTaskInput, coll)
	case 5:
		fmt.Print("Enter index: ")
		scanner.Scan()
		id := scanner.Text()
		deleteTask(id, coll)
	case 6:
		os.Exit(0)
	default:
		fmt.Println("Imvalid choice")
	}
  }
  
}

