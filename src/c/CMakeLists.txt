set(TARGET_NAME whatsmeow)
add_library(${TARGET_NAME} SHARED 
    "gowhatsapp.h"
    "init.c"
    "login.c"
    "qrcode.c"
    "bridge.h"
    "bridge.c"
    "process_message.c"
    "display_message.c"
    "groups.c"
    "blist.c"
    "constants.h"
    "constants.c"
    "send_message.c"
    "handle_attachment.c"
    "send_file.c"
    "presence.c"
    "options.c"
    "receipt.c"
    "pixbuf.c"
    "pixbuf.h"
    "commands.c"
    "opusreader.c"
    "opusreader.h"
)
add_dependencies(${TARGET_NAME} purple-go-whatsapp)

file(READ "${CMAKE_SOURCE_DIR}/VERSION" PLUGIN_VERSION)
target_compile_definitions(${TARGET_NAME} PRIVATE PLUGIN_VERSION=${PLUGIN_VERSION})
target_include_directories(${TARGET_NAME} PRIVATE ${GO_WHATSAPP_INCLUDE_DIRS} ${PURPLE_INCLUDE_DIRS} ${PIXBUF_INCLUDE_DIRS} ${OPUSFILE_INCLUDE_DIRS})
target_link_libraries(${TARGET_NAME} PRIVATE ${GO_WHATSAPP_LIBRARIES} ${PURPLE_LIBRARIES} ${PIXBUF_LIBRARIES} ${OPUSFILE_LIBRARIES})
set_target_properties(${TARGET_NAME} PROPERTIES PREFIX "lib")

if (WIN32)
    install(TARGETS ${TARGET_NAME} RUNTIME DESTINATION "${PURPLE_PLUGIN_DIR}")
else()
    install(TARGETS ${TARGET_NAME} DESTINATION "${PURPLE_PLUGIN_DIR}")
endif()
