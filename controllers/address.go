package controllers

import(

)

func EditWorkAddress() ginHandlerFunc{

}

func DeleteAddress() gin.HandlerFunc{
	//creating an empty object at the ID to be deleted
	return func(c *gin.Context){
		user_id := c.Query("id")

		if user_id == ""{
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error" : "Invalid Search Index"})
			c.Abort()
			return
		}
		addresses := make([]models.Address, 0)
		user_id, err := primitive.ObjectIDFromHex(user_id)
		if err!=nil{
			c.IndentedJSON(500, "Internal Server Error")
		}
		var ctx, cancel = context.WithTimeOut(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D(primitive.E{Key:"_id", Value: usert_id})
		bson.D{Key:"$set", Value: bson.D{primitive.E{Key:"address", Value: addresses}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err!=nil{
			c.IndentedJSON(404, "Wrong command")
			return
		}
		defer cancel()
		ctx.Done()
		c.ItendedJSON(200, "Successfully deleted")
	}
}