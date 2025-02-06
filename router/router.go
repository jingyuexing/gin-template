package router

import "github.com/gin-gonic/gin"

type RouterBoot = func(public *gin.RouterGroup, private *gin.RouterGroup)

type RouterBootList = []RouterBoot
