package services

import (
	"database/sql"
	// "strings"
	"fmt"
	"time"

	"shuttle/models/dto"
	"shuttle/models/entity"
	"shuttle/repositories"

	"github.com/google/uuid"
)

type StudentServiceInterface interface {
	AddPermittedSchoolStudentWithParents(req dto.StudentRequestDTO, username string) error
}

type StudentService struct {
	studentRepository repositories.StudentRepositoryInterface
	userRepository    repositories.UserRepositoryInterface
}

func NewStudentService(studentRepository repositories.StudentRepositoryInterface, userRepository repositories.UserRepositoryInterface) StudentService {
	return StudentService{
		studentRepository: studentRepository,
		userRepository:    userRepository,
	}
}

func (s *StudentService) AddPermittedSchoolStudentWithParents(req dto.AddStudentWithParentRequestDTO, createdBy string) (uuid.UUID, uuid.UUID, error) {
	tx, err := s.studentRepository.BeginTransaction()
	if err != nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("error beginning transaction: %w", err)
	}

	var transactionErr error
	defer func() {
		if transactionErr != nil {
			tx.Rollback()
		} else {
			transactionErr = tx.Commit()
		}
	}()

	// Validate parent data
	parentUUID := uuid.New()
	hashedPassword, err := hashPassword(req.Parent.Password)
	if err != nil {
		transactionErr = fmt.Errorf("error hashing parent password: %w", err)
		return uuid.Nil, uuid.Nil, transactionErr
	}

	parentEntity := entity.User{
		ID:        time.Now().UnixMilli()*1e6 + int64(parentUUID.ID()%1e6),
		UUID:      parentUUID,
		Username:  req.Parent.Username, // Assuming email is used as username
		Email:     req.Parent.Email,
		Password:  hashedPassword,
		Role:      entity.Parent,
		RoleCode:  "P",
		CreatedBy: sql.NullString{String: createdBy, Valid: createdBy != ""},
		Details: map[string]interface{}{
			"first_name": req.Parent.FirstName,
			"last_name":  req.Parent.LastName,
			"phone":      req.Parent.Phone,
			"address":    req.Parent.Address,
		},
	}

	parentUUID, err = s.userRepository.SaveUser(tx, parentEntity)
	if err != nil {
		transactionErr = fmt.Errorf("error saving parent: %w", err)
		return uuid.Nil, uuid.Nil, transactionErr
	}

	// Validate student data
	studentUUID := uuid.New()
	studentEntity := entity.Student{
		ID:         time.Now().UnixMilli()*1e6 + int64(studentUUID.ID()%1e6),
		UUID:       studentUUID,
		FirstName:  req.Student.FirstName,
		LastName:   req.Student.LastName,
		Grade:      req.Student.Grade,
		Gender:     req.Student.Gender,	
		ParentUUID: sql.NullString{String: parentUUID.String(), Valid: true},
		SchoolUUID: uuid.MustParse(req.Student.SchoolUUID),
		CreatedBy:  sql.NullString{String: createdBy, Valid: createdBy != ""},
	}

	studentUUID, err = s.studentRepository.SaveStudent(tx, studentEntity)
	if err != nil {
		transactionErr = fmt.Errorf("error saving student: %w", err)
		return uuid.Nil, uuid.Nil, transactionErr
	}

	return studentUUID, parentUUID, nil
}

// package services

// import (
// 	"context"

// 	// "shuttle/errors"
// 	"shuttle/databases"
// 	"shuttle/models"

// 	"github.com/spf13/viper"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// func GetAllPermitedSchoolStudentsWithParents(schoolID primitive.ObjectID) ([]models.SchoolStudentParentResponse, error) {
// 	client, err := databases.MongoConnection()
// 	if err != nil {
// 		return nil, err
// 	}

// 	collection := client.Database(viper.GetString("MONGO_DB")).Collection("students")

// 	var students []models.Student
// 	cursor, err := collection.Find(context.Background(), bson.M{"school_id": schoolID})
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(context.Background())

// 	for cursor.Next(context.Background()) {
// 		var student models.Student
// 		if err := cursor.Decode(&student); err != nil {

// 			return nil, err
// 		}
// 		students = append(students, student)
// 	}

// 	if err := cursor.Err(); err != nil {
// 		return nil, err
// 	}

// 	var Parents []models.ParentResponse
// 	parentCollection := client.Database(viper.GetString("MONGO_DB")).Collection("users")
// 	for _, student := range students {
// 		var parent models.ParentResponse
// 		err := parentCollection.FindOne(context.Background(), bson.M{"_id": student.ParentID}, options.FindOne().SetProjection(bson.M{"password": 0})).Decode(&parent)
// 		if err != nil {

// 			return nil, err
// 		}
// 		Parents = append(Parents, parent)
// 	}

// 	var schoolStudents []models.SchoolStudentParentResponse
// 	for i, student := range students {
// 		schoolStudents = append(schoolStudents, models.SchoolStudentParentResponse{
// 			Student: student,
// 			Parent:  Parents[i],
// 		})
// 	}

// 	return schoolStudents, nil
// }

// func AddPermittedSchoolStudentWithParents(student models.SchoolStudentRequest, schoolID primitive.ObjectID, username string) error {
// 	client, err := databases.MongoConnection()
// 	if err != nil {
// 		return err
// 	}

// 	collection := client.Database(viper.GetString("MONGO_DB")).Collection("users")
// 	var existingParent models.User
// 	// Email as unique identifier for parent
// 	err = collection.FindOne(context.Background(), bson.M{"email": student.Parent.Email, "role": models.Parent}).Decode(&existingParent)

// 	// If parent is not yet added, add the parent
// 	var parentID primitive.ObjectID
// 	if err == nil {
// 		parentID = existingParent.ID
// 	} else if err == mongo.ErrNoDocuments {
// 		parentUser := student.Parent
// 		parentUser.Role = models.Parent

// 		parentUser.Details = &models.ParentDetails{
// 			Children: []primitive.ObjectID{},
// 		}

// 		parentID, err = AddUser(parentUser, username)
// 		if err != nil {

// 			return err
// 		}
// 	} else {
// 		return err
// 	}

// 	studentDocument := bson.D{
// 		{Key: "first_name", Value: student.Student.FirstName},
// 		{Key: "last_name", Value: student.Student.LastName},
// 		{Key: "class", Value: student.Student.Class},
// 		{Key: "parent_id", Value: parentID},
// 		{Key: "school_id", Value: schoolID},
// 	}

// 	studentsCollection := client.Database(viper.GetString("MONGO_DB")).Collection("students")
// 	result, err := studentsCollection.InsertOne(context.Background(), studentDocument)
// 	if err != nil {
// 		return err
// 	}

// 	studentID := result.InsertedID.(primitive.ObjectID)

// 	_, err = collection.UpdateOne(
// 		context.Background(),
// 		bson.M{"_id": parentID},
// 		bson.M{"$push": bson.M{"details.children_id": studentID}},
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func UpdatePermittedSchoolStudentWithParents(id string, student models.SchoolStudentRequest, schoolID primitive.ObjectID) error {
//     client, err := databases.MongoConnection()
//     if err != nil {
//         return err
//     }

//     collection := client.Database(viper.GetString("MONGO_DB")).Collection("students")
//     objectID, err := primitive.ObjectIDFromHex(id)
//     if err != nil {
//         return err
//     }

// 	if err := CheckStudentAvailability(objectID, schoolID); err != nil {
//         return err
//     }

// 	// Pipeline to get the student and parent details
//     var existingStudent models.SchoolStudentRequest
//     pipeline := mongo.Pipeline{
//         bson.D{{Key: "$match", Value: bson.M{"_id": objectID, "school_id": schoolID}}},
//         bson.D{{Key: "$lookup", Value: bson.M{
//             "from":         "users",
//             "localField":   "parent_id",
//             "foreignField": "_id",
//             "as":           "parent",
//         }}},
//         bson.D{{Key: "$unwind", Value: bson.M{"path": "$parent"}}},
//     }

// 	// Aggregate the pipeline
//     cursor, err := collection.Aggregate(context.Background(), pipeline)
//     if err != nil {
//         return err
//     }
//     defer cursor.Close(context.Background())

//     if cursor.Next(context.Background()) {
//         if err := cursor.Decode(&existingStudent); err != nil {

//             return err
//         }
//     }

//     if (models.User{}) == existingStudent.Parent {
//         return errors.New("parent details are not available", 404)
//     }

//     updateStudent := bson.M{
//         "first_name": student.FirstName,
//         "last_name":  student.LastName,
//     }

//     _, err = collection.UpdateOne(context.Background(), bson.M{"_id": objectID, "school_id": schoolID}, bson.M{"$set": updateStudent})
//     if err != nil {
//         return err
//     }

// 	// If the parent details are changed, update the parent details
//     if (models.User{}) != student.Parent && student.Parent != existingStudent.Parent {
//         parentCollection := client.Database(viper.GetString("MONGO_DB")).Collection("users")
//         parentID := existingStudent.Parent.ID

//         updateParent := bson.M{
//             "first_name": student.Parent.FirstName,
//             "last_name":  student.Parent.LastName,
//             "email":      student.Parent.Email,
// 			"phone":    student.Parent.Phone,
//             "address":  student.Parent.Address,
//             "details": bson.M{
//                 "children": []primitive.ObjectID{objectID},
//             },
//         }

//         _, err = parentCollection.UpdateOne(context.Background(), bson.M{"_id": parentID}, bson.M{"$set": updateParent})
//         if err != nil {

//             return err
//         }
//     } else { // Else, parent remains the same
//         _, err = client.Database(viper.GetString("MONGO_DB")).Collection("users").UpdateOne(
//             context.Background(),
//             bson.M{"_id": existingStudent.Parent.ID},
//             bson.M{"$addToSet": bson.M{"parent_details.children": objectID}},
//         )
//         if err != nil {

//             return err
//         }
//     }

//     return nil
// }

// func DeletePermittedSchoolStudentWithParents(id string, schoolID primitive.ObjectID) error {
// 	client, err := databases.MongoConnection()
// 	if err != nil {
// 		return err
// 	}

// 	collection := client.Database(viper.GetString("MONGO_DB")).Collection("students")
// 	objectID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}

// 	if err := CheckStudentAvailability(objectID, schoolID); err != nil {
// 		return err
// 	}

// 	var student models.Student
// 	err = collection.FindOne(context.Background(), bson.M{"_id": objectID, "school_id": schoolID}).Decode(&student)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objectID, "school_id": schoolID})
// 	if err != nil {
// 		return err
// 	}

// 	// Remove the student from the parent's children list
// 	parentCollection := client.Database(viper.GetString("MONGO_DB")).Collection("users")
// 	_, err = parentCollection.UpdateOne(
// 		context.Background(),
// 		bson.M{"_id": student.ParentID},
// 		bson.M{"$pull": bson.M{"parent_details.children": objectID}},
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	var parent models.User
// 	err = parentCollection.FindOne(context.Background(), bson.M{"_id": student.ParentID}).Decode(&parent)
// 	if err != nil {
// 		return err
// 	}

// 	// If the children array is empty, delete the parent
// 	if len(parent.Details.(models.ParentDetails).Children) == 0 {
// 		_, err = parentCollection.DeleteOne(context.Background(), bson.M{"_id": student.ParentID})
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func CheckPermittedSchoolAccess(userID string) (primitive.ObjectID, error) {
// 	client, err := databases.MongoConnection()
// 	if err != nil {
// 		return primitive.NilObjectID, err
// 	}

// 	objectID, err := primitive.ObjectIDFromHex(userID)
// 	if err != nil {
// 		return primitive.NilObjectID, err
// 	}

// 	collection := client.Database(viper.GetString("MONGO_DB")).Collection("users")

// 	var user models.User
// 	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
// 	if err != nil {
// 		return primitive.NilObjectID, err
// 	}

// 	var schoolAdminDetails models.SchoolAdminDetails
// 	detailsBytes, err := bson.Marshal(user.Details)
// 	if err != nil {
// 		return primitive.NilObjectID, err
// 	}

// 	err = bson.Unmarshal(detailsBytes, &schoolAdminDetails)
// 	if err != nil || schoolAdminDetails.SchoolID.IsZero() {
// 		return primitive.NilObjectID, err
// 	}

// 	return schoolAdminDetails.SchoolID, nil
// }

// func CheckStudentAvailability(studentID primitive.ObjectID, schoolID primitive.ObjectID) error {
//     client, err := databases.MongoConnection()
//     if err != nil {
//         return err
//     }

//     collection := client.Database(viper.GetString("MONGO_DB")).Collection("students")

//     var student models.SchoolStudentRequest
//     err = collection.FindOne(context.Background(), bson.M{"_id": studentID, "school_id": schoolID}).Decode(&student)
//     if err != nil {
//         if err == mongo.ErrNoDocuments {
//             return err
// 		}
//         return err
//     }

//     return nil
// }
