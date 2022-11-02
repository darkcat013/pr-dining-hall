package config

import "time"

const KITCHEN_URL = "http://localhost:8080/order"

//const KITCHEN_URL = "http://host.docker.internal:8080/order"

const FOOD_ORDERING_SERVICE_URL = "http://localhost:8088"

// const FOOD_ORDERING_SERVICE_URL = "http://host.docker.internal:8088"

const PORT = "8081"
const LOGS_ENABLED = true

const TIME_UNIT = time.Millisecond * TIME_UNIT_COEFF
const TIME_UNIT_COEFF = 100

const TABLES = 10
const WAITERS = 4

const MENU_PATH = "config/menu.json"

const MAX_FOODS = 10
const MAX_PREP_TIME_COEFF = 1.3

const TABLE_OCCUPATION_TIME_MIN = 15
const TABLE_OCCUPATION_TIME_MAX = 30
const TABLE_ORDERING_TIME_MIN = 5
const TABLE_ORDERING_TIME_MAX = 10
const TABLE_PICKING_ORDER_TIME = 3

const WAITER_TAKING_ORDER_TIME_MIN = 2
const WAITER_TAKING_ORDER_TIME_MAX = 4

const RESTAURANT_NAME = "Gamer food"
const RESTAURANT_ID = 1

const ADDRESS = "http://localhost:8081"

// const ADDRESS = "http://host.docker.internal:8081"
