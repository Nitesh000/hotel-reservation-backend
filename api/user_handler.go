package api

import (
	"errors"

	"github.com/Nitesh000/hotel-reservation-backend/db"
	"github.com/Nitesh000/hotel-reservation-backend/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandePutUser(c *fiber.Ctx) error {
	return nil
}

func (h *UserHandler) HandeDeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	if err := h.userStore.DeleteUser(c.Context(), userId); err != nil {
		return err
	}
	return c.JSON(map[string]string{"message": "user has been deleted!"})
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.userStore.GetUserById(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(404).JSON(map[string]string{"msg": "user not found with id " + id + "."})
		}
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserPrams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}

	insertedUser, err := h.userStore.PostUser(c.Context(), user)
	if err != nil {
		return err
	}

	return c.JSON(&insertedUser)
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	// var update bson.M
	var userData types.UpdateUserParams
	userId := c.Params("id")

	// NOTE: converting the string id to object id
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return c.JSON(map[string]string{"error": "wrong account Id."})
	}

	if err := c.BodyParser(&userData); err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	if err := h.userStore.UpdateUser(c.Context(), filter, userData); err != nil {
		return err
	}
	return c.JSON(map[string]string{"updated user id": userId})
}
