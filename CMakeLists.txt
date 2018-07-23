cmake_minimum_required( VERSION 2.8 FATAL_ERROR )

set( CMAKE_VERBOSE_MAKEFILE ON )

project( SparkTest )

get_filename_component( CINDER_PATH "/Users/henrique/Documents/Dev/cpp/Cinder" ABSOLUTE )
get_filename_component( APP_PATH "${CMAKE_CURRENT_SOURCE_DIR}/../../" ABSOLUTE )

include( "${CINDER_PATH}/proj/cmake/modules/cinderMakeApp.cmake" )

ci_make_app(

	SOURCES     ${APP_PATH}/src/SparkTestApp.cpp
	CINDER_PATH ${CINDER_PATH}
	BLOCKS		OSC
)
