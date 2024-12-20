package handlers

import (
	"RMS_deploy/initializers"
	Log "RMS_deploy/log"
	"RMS_deploy/models"
	"database/sql"
	"encoding/json"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func UserExist(c *gin.Context, email string) (user models.Users, err error) {
	log := Log.GetLogger(c)
	query := "SELECT id,email,password_hash FROM users WHERE email = $1"

	err = initializers.DB.GetContext(c, &user, query, email)
	if err != nil {
		log.Warn("Executing SQL query",
			zap.String("query", query),
			zap.String("userID", email),
		)

	}

	return user, err

}

func CreateUser(c *gin.Context, email string, hash string, name string, role string, userCreateId uint) error {
	log := Log.GetLogger(c)
	var user models.Users
	query := `INSERT INTO users ( email, password_hash,created_at,updated_at,name,role,created_by)
			  VALUES (TRIM($1), TRIM($2),$3,$4,$5,$6,$7) RETURNING id`

	err := initializers.DB.Get(&user.ID, query, email, hash, time.Now(), time.Now(), name, role, userCreateId)
	if err != nil {
		log.Warn("Executing SQL query",
			zap.String("query", query),
			zap.String("userID", email),
			zap.String("hash", hash),
			zap.String("role", role),
			zap.String("name", name),
			zap.String("created_by", role),
		)
	}
	return err
}

func CreateRestaurant(c *gin.Context, name string, longitude float64, latitude float64, userId uint) (restaurant models.Restaurant, err error) {
	log := Log.GetLogger(c)

	//query := `INSERT INTO restaurants ( name,longitude,latitude, owner_id,created_at,updated_at)
	//		  VALUES ($1, $2,$3,$4,$5,$6) RETURNING id`
	query := `INSERT INTO restaurants ( name,longitude,latitude, owner_id,created_at,updated_at)
           VALUES ($1, $2,$3,$4,$5,$6) RETURNING id, name, longitude, latitude, owner_id, created_at, updated_at`

	err = initializers.DB.GetContext(c, &restaurant, query, name, longitude, latitude, userId, time.Now(), time.Now())
	if err != nil {
		log.Warn("Executing SQL query",
			zap.String("query", query),
			zap.String("name", name),
			zap.Float64("longitude", longitude),
			zap.Float64("latitude", latitude),
			zap.Uint("userId", userId),
		)
	}
	return restaurant, err
}

func CreateDish(c *gin.Context, name string, price decimal.Decimal, restaurantId uint, ownerId uint, tags []string) (dish models.Dish, err error) {
	log := Log.GetLogger(c)
	query := `INSERT INTO dishes (name, price, restaurant_id, owner_id, tags, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, name, price, restaurant_id, owner_id, tags, created_at, updated_at`

	err = initializers.DB.GetContext(c, &dish, query, name, price, restaurantId, ownerId, pq.Array(tags), time.Now(), time.Now())
	if err != nil {
		log.Warn("Executing SQL query",
			zap.String("query", query),
			zap.String("name", name),
			zap.Any("price", price),
			zap.Uint("restaurantId", restaurantId),
			zap.Uint("ownerId", ownerId),
			zap.Any("tags", tags),
		)
	}
	return dish, err
}

//func GetAll(c *gin.Context, userId uint) ([]models.Note, error) {
//	log := Log.GetLogger(c)
//	var notes []models.Note
//	query := `SELECT id,title, body,created_at,updated_at, user_id
//        FROM notes
//        WHERE user_id = $1 and deleted_at is null`
//	err := initializers.DB.Select(&notes, query, userId)
//	if err != nil {
//		log.Warn(err.Error(),
//			zap.String("query", query),
//			zap.Uint("userId", userId),
//		)
//	}
//
//	return notes, err
//}

//func GetOne(c *gin.Context, id string) (models.Note, error) {
//	log := Log.GetLogger(c)
//	var note models.Note
//	query := "SELECT id,title,body, created_at,updated_at,deleted_at,user_id FROM notes WHERE id = $1 and deleted_at is null"
//
//	err := initializers.DB.GetContext(c, &note, query, id)
//	fmt.Println(err)
//	if err != nil {
//		log.Warn("Executing SQL query",
//			zap.String("query", query),
//			zap.String("id", id),
//		)
//	}
//	return note, err
//}

func UpdateOne(c *gin.Context, title string, body string, id string) (sql.Result, error) {
	log := Log.GetLogger(c)
	updateQuery := `
        UPDATE notes 
        SET title = $1,body = $2
        WHERE id = $3`
	result, err := initializers.DB.Exec(updateQuery, title, body, id)
	if err != nil {
		log.Warn("Executing SQL query",
			zap.String("query", updateQuery),
			zap.String("title", title),
			zap.String("body", body),
			zap.String("id", id),
		)
	}
	return result, err
}

func DeleteOne(c *gin.Context, id string) (sql.Result, error) {
	log := Log.GetLogger(c)
	query := `
    UPDATE notes 
    SET deleted_at = $1 
    WHERE id = $2 `
	result, err := initializers.DB.Exec(query, time.Now(), id)
	if err != nil {
		log.Warn("Executing SQL query",
			zap.String("id", id),
			zap.String("query", query),
		)
	}
	return result, err
}

//func ExcelRead(c *gin.Context, notes []models.Note) error {
//	log := Log.GetLogger(c)
//	tx, err := initializers.DB.Beginx()
//	if err != nil {
//		zap.L().Info(err.Error())
//		return err
//	}
//
//	for _, note := range notes {
//		query := `INSERT INTO notes (title, body,user_id,created_at,updated_at) VALUES (:title, :body,:user_id,:created_at,:updated_at)`
//		_, err = tx.NamedExec(query, note)
//		if err != nil {
//			tx.Rollback()
//			log.Warn(err.Error(),
//				zap.String("query", query),
//				zap.String("title", note.Title),
//				zap.String("body", note.Body),
//				zap.Uint("userId", note.UserID),
//			)
//			return err
//		}
//	}
//	tx.Commit()
//	return err
//}

func AddAddress(c *gin.Context, name string, longitude float64, latitude float64, userId uint) (err error) {
	log := Log.GetLogger(c)

	var addr models.Address
	query := `INSERT INTO addresses (name, longitude, latitude, user_id, created_at)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`

	// Execute the query and get the inserted ID
	err = initializers.DB.GetContext(c, &addr, query, name, longitude, latitude, userId, time.Now())
	if err != nil {
		log.Error("Failed to execute SQL query",
			zap.String("query", query),
			zap.String("name", name),
			zap.Float64("longitude", longitude),
			zap.Float64("latitude", latitude),
			zap.Uint32("userId", uint32(userId)),
			zap.Error(err),
		)
		return err
	}

	log.Info("Address added successfully",
		zap.Uint("address_id", addr.ID),
		zap.String("name", name),
	)
	return nil

}

func GetAllRest(c *gin.Context) ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	query := `SELECT 
    id,
    name,
    owner_id,
    latitude,
    longitude,
    created_at,
    updated_at
FROM restaurants
WHERE deleted_at IS NULL;`
	err := initializers.DB.SelectContext(c, &restaurants, query)
	if err != nil {
		zap.L().Warn(err.Error(),
			zap.String("query", query),
		)
	}

	return restaurants, err
}

func GetAllDish(c *gin.Context, id string) (dishes []models.Dish, err error) {
	query := `SELECT id,name,price, restaurant_id,owner_id,tags,created_at,updated_at
				FROM dishes
				WHERE  restaurant_id= $1 and deleted_at IS NULL
				ORDER BY id;`
	err = initializers.DB.SelectContext(c, &dishes, query, id)
	if err != nil {
		zap.L().Warn(err.Error(),
			zap.String("query", query),
		)
	}
	return dishes, err
}
func GetRestById(c *gin.Context, id string) (Rest models.Restaurant, err error) {
	log := Log.GetLogger(c)

	query := "SELECT id,\n    name,\n    owner_id,\n    latitude,\n    longitude,\n    created_at,\n    updated_at\nFROM restaurants WHERE id = $1 and deleted_at is null"

	err = initializers.DB.GetContext(c, &Rest, query, id)
	if err != nil {
		log.Warn(err.Error(),
			zap.String("query", query),
			zap.String("id", id),
		)
	}
	return Rest, err
}

func GetAllAddr(c *gin.Context, id string) (addrs []models.Address, err error) {

	query := `SELECT id,name,user_id, latitude,longitude,created_at
				FROM addresses
				WHERE  user_id= $1 and deleted_at IS NULL
				ORDER BY id;`
	err = initializers.DB.SelectContext(c, &addrs, query, id)
	if err != nil {
		zap.L().Warn(err.Error(),
			zap.String("query", query),
		)
	}

	return addrs, err
}

func GetAllSubAdmins(c *gin.Context) (user []models.Users, err error) {
	query := `SELECT 
    id,
    name,
    email,
    role,
    created_by,
    created_at,
    updated_at
FROM users
WHERE deleted_at IS NULL and role='sub-admin';`
	err = initializers.DB.SelectContext(c, &user, query)
	if err != nil {
		zap.L().Warn(err.Error(),
			zap.String("query", query),
		)
	}
	return user, err

}

func GetAllRestaurantAdmin(c *gin.Context) (restaurants []models.Restaurant, err error) {

	//log := Log.GetLogger(c)
	// SQL query to fetch restaurants and their dishes as JSON
	//query := `
	//SELECT
	//	r.id AS id,
	//	r.name AS name,
	//	r.latitude AS latitude,
	//	r.longitude AS longitude,
	//	r.owner_id AS owner_id,
	//	r.created_at AS created_at,
	//	r.updated_at AS updated_at,
	//	COALESCE(
	//		JSON_AGG(
	//			JSON_BUILD_OBJECT(
	//				'id', d.id,
	//				'name', d.name,
	//				'price', d.price,
	//				'tags', d.tags,
	//				'restaurant_id', d.restaurant_id,
	//				'owner_id', d.owner_id,
	//				'created_at', d.created_at,
	//				'updated_at', d.updated_at
	//			)
	//		) FILTER (WHERE d.id IS NOT NULL), '[]'
	//	) AS dishes
	//FROM
	//	restaurants r
	//LEFT JOIN
	//	dishes d ON r.id = d.restaurant_id AND d.deleted_at IS NULL
	//WHERE
	//	r.deleted_at IS NULL
	//GROUP BY
	//	r.id
	//ORDER BY
	//	r.id;
	//`
	//
	//// Execute the query and scan results into the restaurants slice
	//err = initializers.DB.SelectContext(c, &restaurants, query)
	//if err != nil {
	//	log.Warn(err.Error(),
	//		zap.String("query", query),
	//	)
	//}
	//return restaurants, nil

	query := `
	SELECT 
		r.id AS id,
		r.name AS name,
		r.latitude AS latitude,
		r.longitude AS longitude,
		r.owner_id AS owner_id,
		r.created_at AS created_at,
		r.updated_at AS updated_at,
		COALESCE(
			JSON_AGG(
				JSON_BUILD_OBJECT(
					'id', d.id,
					'name', d.name,
					'price', d.price,
					'tags', d.tags,
					'restaurant_id', d.restaurant_id,
					'owner_id', d.owner_id,
					'created_at', d.created_at,
					'updated_at', d.updated_at
				)
			) FILTER (WHERE d.id IS NOT NULL), '[]'
		) AS dishes
	FROM 
		restaurants r
	LEFT JOIN 
		dishes d ON r.id = d.restaurant_id AND d.deleted_at IS NULL
	WHERE 
		r.deleted_at IS NULL
	GROUP BY 
		r.id
	ORDER BY 
		r.id;
	`
	rows, err := initializers.DB.QueryContext(c, query)
	if err != nil {
		zap.L().Error("Error executing query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	// Process rows manually
	for rows.Next() {
		var restaurant models.Restaurant
		var dishesJSON []byte

		// Scan the restaurant fields and the JSON-encoded dishes
		err := rows.Scan(
			&restaurant.ID,
			&restaurant.Name,
			&restaurant.Latitude,
			&restaurant.Longitude,
			&restaurant.OwnerID,
			&restaurant.CreatedAt,
			&restaurant.UpdatedAt,
			&dishesJSON, // Raw JSON for dishes
		)
		if err != nil {
			zap.L().Error("Error scanning row", zap.Error(err))
			return nil, err
		}

		// Decode the JSON into the `Dishes` field
		err = json.Unmarshal(dishesJSON, &restaurant.Dishes)
		if err != nil {
			zap.L().Error("Error unmarshalling dishes JSON", zap.Error(err))
			return nil, err
		}

		restaurants = append(restaurants, restaurant)
	}

	if err = rows.Err(); err != nil {
		zap.L().Error("Error in rows iteration", zap.Error(err))
		return nil, err
	}

	return restaurants, nil
}

// func GetAllRestaurantSubAdmin(c *gin.Context, userId string) (restaurants []models.Restaurant, err error) {
//
//		log := Log.GetLogger(c)
//		// SQL query to fetch restaurants and their dishes as JSON
//		query := `
//		SELECT
//			r.id AS id,
//			r.name AS name,
//			r.latitude AS latitude,
//			r.longitude AS longitude,
//			r.owner_id AS owner_id,
//			r.created_at AS created_at,
//			r.updated_at AS updated_at,
//			COALESCE(
//				JSON_AGG(
//					JSON_BUILD_OBJECT(
//						'id', d.id,
//						'name', d.name,
//						'price', d.price,
//						'tags', d.tags,
//						'restaurant_id', d.restaurant_id,
//						'owner_id', d.owner_id,
//						'created_at', d.created_at,
//						'updated_at', d.updated_at
//					)
//				) FILTER (WHERE d.id IS NOT NULL), '[]'
//			) AS dishes
//		FROM
//			restaurants r
//		LEFT JOIN
//			dishes d ON r.id = d.restaurant_id AND d.deleted_at IS NULL
//		WHERE
//			r.deleted_at IS NULL AND r.owner_id = $1
//		GROUP BY
//			r.id
//		ORDER BY
//			r.id;
//		`
//
//		// Execute the query and scan results into the restaurants slice
//		err = initializers.DB.SelectContext(c, &restaurants, query, userId)
//		if err != nil {
//			log.Warn(err.Error(),
//				zap.String("query", query),
//			)
//		}
//		return restaurants, nil
//	}
func GetAllRestaurantSubAdmin(c *gin.Context, userId string) ([]models.Restaurant, error) {
	log := Log.GetLogger(c)
	var restaurants []models.Restaurant

	// Step 1: Fetch all restaurants for the given user
	queryRestaurants := `
		SELECT
			id, name, latitude, longitude, owner_id, created_at, updated_at
		FROM
			restaurants
		WHERE
			deleted_at IS NULL AND owner_id = $1
		ORDER BY
			id;
	`

	err := initializers.DB.SelectContext(c, &restaurants, queryRestaurants, userId)
	if err != nil {
		log.Warn("Failed to fetch restaurants",
			zap.String("query", queryRestaurants),
			zap.Error(err),
		)
		return nil, err
	}

	// If no restaurants are found, return an empty slice
	if len(restaurants) == 0 {
		return restaurants, nil
	}

	// Step 2: Extract restaurant IDs
	var restaurantIDs []uint
	for _, restaurant := range restaurants {
		restaurantIDs = append(restaurantIDs, restaurant.ID)
	}

	// Step 3: Fetch dishes for the restaurants
	queryDishes := `
		SELECT
			id, name, price, tags, restaurant_id, owner_id, created_at, updated_at
		FROM
			dishes
		WHERE
			deleted_at IS NULL AND restaurant_id = ANY($1)
	`

	var dishes []models.Dish
	err = initializers.DB.SelectContext(c, &dishes, queryDishes, pq.Array(restaurantIDs))
	if err != nil {
		log.Warn("Failed to fetch dishes",
			zap.String("query", queryDishes),
			zap.Error(err),
		)
		return nil, err
	}

	// Step 4: Map dishes to their respective restaurants
	dishMap := make(map[uint][]models.Dish)
	for _, dish := range dishes {
		dishMap[dish.RestaurantID] = append(dishMap[dish.RestaurantID], dish)
	}

	// Step 5: Assign dishes to restaurants
	for i := range restaurants {
		restaurants[i].Dishes = dishMap[restaurants[i].ID]
	}

	return restaurants, nil
}
func GetAllUsersAdmin(c *gin.Context) (users []models.Users, err error) {

	log := Log.GetLogger(c)
	query := `SELECT 
    id,
    name,
    email,
    role,
    created_by,
    created_at,
    updated_at
FROM users
WHERE deleted_at IS NULL and role='users';`
	err = initializers.DB.SelectContext(c, &users, query)
	if err != nil {
		log.Warn(err.Error(),
			zap.String("query", query),
		)
	}
	return users, err
}

func GetAllUsersSubAdmin(c *gin.Context, userId string) (users []models.Users, err error) {

	log := Log.GetLogger(c)
	query := `SELECT 
    id,
    name,
    email,
    role,
    created_by,
    created_at,
    updated_at
FROM users
WHERE deleted_at IS NULL and role='user' and created_by=$1;`
	err = initializers.DB.SelectContext(c, &users, query, userId)
	if err != nil {
		log.Warn(err.Error(),
			zap.String("query", query),
		)
	}
	return users, err
}
