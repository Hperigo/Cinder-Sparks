cmake_minimum_required( VERSION 2.8 FATAL_ERROR )
set( CMAKE_VERBOSE_MAKEFILE ON )

project( __PROJECT_NAME__ )

get_filename_component( CINDER_PATH __CINDER_PATH__ ABSOLUTE )
get_filename_component( APP_PATH "${CMAKE_CURRENT_SOURCE_DIR}" ABSOLUTE )

include( "${CINDER_PATH}/proj/cmake/modules/cinderMakeApp.cmake" )

ci_make_app(
	SOURCES     ${APP_PATH}/src/__PROJECT_NAME__App.cpp
	CINDER_PATH ${CINDER_PATH}
)
