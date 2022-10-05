package config

import "time"

const KITCHEN_URL = "http://localhost:8080/order"
const LOGS_ENABLED = true

const TIME_UNIT = time.Millisecond * TIME_UNIT_COEFF
const TIME_UNIT_COEFF = 1

const TABLES = 10
const WAITERS = 4

const MENU_PATH = "config/menu.json"

const MAX_FOODS = 10
const MAX_PREP_TIME_COEFF = 1.3

const TABLE_OCCUPATION_TIME_MIN = 0
const TABLE_OCCUPATION_TIME_MAX = 1
const TABLE_ORDERING_TIME_MIN = 0
const TABLE_ORDERING_TIME_MAX = 1
const TABLE_PICKING_ORDER_TIME = 1

const WAITER_TAKING_ORDER_TIME_MIN = 2
const WAITER_TAKING_ORDER_TIME_MAX = 4
