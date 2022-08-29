package main

type credentials struct {
	Admin    string `env:"PALADA_ADMIN" envDefault:"admin"`
	Password string `env:"PALADA_PASSWORD" envDefault:"password"`
	Database string `env:"PALADA_DATABASE" envDefault:"timetable"`
}
