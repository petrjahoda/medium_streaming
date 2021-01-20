package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type Data struct {
	gorm.Model
	Data string
}

func checkDatabase() {
	databaseCheck := false
	for !databaseCheck {
		productionDatabase, err := gorm.Open(postgres.Open(databaseConnection), &gorm.Config{})
		productionDB, _ := productionDatabase.DB()
		if err != nil {
			fmt.Println("Problem opening database, looks like it does not exist")
			database, err := gorm.Open(postgres.Open(connection), &gorm.Config{})
			sqlDB, _ := database.DB()
			if err != nil {
				fmt.Println("Problem opening main  database: " + err.Error())
				time.Sleep(1 * time.Second)
				continue
			}
			fmt.Println("Creating medium database")
			database.Exec("CREATE DATABASE medium;")
			_ = sqlDB.Close()
			continue
		} else {
			fmt.Println("Medium database already exists")
			if !productionDatabase.Migrator().HasTable(&Data{}) {
				fmt.Println("Creating table Data")
				err := productionDatabase.Migrator().CreateTable(&Data{})
				if err != nil {
					fmt.Println("Cannot create table: " + err.Error())
					return
				}
			} else {
				fmt.Println("Updating table Data")
				err := productionDatabase.Migrator().AutoMigrate(&Data{})
				if err != nil {
					fmt.Println("Cannot update table: " + err.Error())
					return
				}
			}
			productionDatabase.Exec("CREATE OR REPLACE FUNCTION notify_event() RETURNS TRIGGER AS\n$$\nDECLARE\n    data json; notification json;\nBEGIN\n    IF (TG_OP = 'DELETE') THEN\n        data = row_to_json(OLD);\n    ELSE\n        data = row_to_json(NEW);\n    END IF;\n    notification = json_build_object('data', data);\n    PERFORM pg_notify('events', notification::text);\n    RETURN NULL;\nEND;\n$$ LANGUAGE plpgsql;")
			productionDatabase.Exec("CREATE TRIGGER products_notify_event\n    AFTER INSERT OR UPDATE OR DELETE\n    ON data\n    FOR EACH ROW\nEXECUTE PROCEDURE notify_event();")
		}
		_ = productionDB.Close()
		databaseCheck = true
		time.Sleep(1 * time.Second)
	}
	fmt.Println("Checking database done")
}
