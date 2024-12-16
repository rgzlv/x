#define _POSIX_C_SOURCE 200112L
#include <stdio.h>
#include <stdbool.h>
#include <stdint.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <time.h>
#include <sys/mman.h>
#include <fcntl.h>

#include <wayland-client.h>
#include "xdg-shell.h"

int debug_client = 1;

typedef struct {
    struct wl_display *display;
    struct wl_registry *registry;
    struct wl_compositor *compositor;
    struct wl_surface *surface;
    struct wl_shm *shm;
    struct wl_shm_pool *shm_pool;

    const struct wl_registry_listener registry_listener;

    struct xdg_wm_base *xdg_wm_base;
    struct xdg_surface *xdg_surface;
    struct xdg_toplevel *xdg_toplevel;

    const struct xdg_surface_listener xdg_surface_listener;
    const struct xdg_toplevel_listener xdg_toplevel_listener;

    char *shmName;
    int shmFd;
    int32_t width;
    int32_t height;
    bool closed;
} client_state;

char *wlc_rndShm(void);
void wlc_registry_global(void *data, struct wl_registry *registry, uint32_t id, const char *interface, uint32_t version);
void wlc_registry_global_remove(void *data, struct wl_registry *registry, uint32_t id);
void wlc_xdg_surface_configure(void *data, struct xdg_surface *xdg_surface, uint32_t serial);
struct wl_buffer *wlc_draw_frame(client_state *cstate);

char *wlc_rndShm(void) {
    char *name = (char *)malloc(18);
    strcpy(name, "/wlc_shm-XXXXXXXX");

    struct timespec ts;
    if (clock_gettime(CLOCK_REALTIME, &ts) != 0) {
        return NULL;
    }

    srand((unsigned int)ts.tv_nsec);
    uint32_t rnd = (uint32_t)rand();
    // Mask with 26 LSB set to 1
    // Bitwise and will make it 8 digits or less so it fits in rnd_str
    uint32_t mask = 67108863;
    rnd = rnd & mask;

    char *rndStr = (char *)malloc(9);
    if (rndStr == NULL) {
        return NULL;
    }
    memset(rndStr, 0, 9);
    if (sprintf(rndStr, "%d", rnd) < 0) {
        return NULL;
    }
 
    int offset = 0;
    for (char *i = &name[9]; *i != '\0'; i++) {
        *i = *(rndStr + offset);
        offset++;
    }
 
    return name;
}

void wlc_registry_global(void *data, struct wl_registry *registry, uint32_t id, const char *interface, uint32_t version) {
    client_state *cstate = (client_state *)data;

    printf("INFO: wlc_registry_global() interface: %s, id: %d\n", interface, id);

    if (strcmp(interface, "wl_compositor") == 0) {
        cstate->compositor = (struct wl_compositor *)wl_registry_bind(registry, id, &wl_compositor_interface, 5);
        return;
    }

    if (strcmp(interface, "xdg_wm_base") == 0) {
        cstate->xdg_wm_base = (struct xdg_wm_base *)wl_registry_bind(registry, id, &xdg_wm_base_interface, 6);
        return;
    }

    if (strcmp(interface, "wl_shm") == 0) {
        cstate->shm = (struct wl_shm *)wl_registry_bind(registry, id, &wl_shm_interface, 1);
    }
}

void wlc_registry_global_remove(void *data, struct wl_registry *registry, uint32_t id)
{
    printf("registry_del_handler() id: %d\n", id);
}

void wlc_xdg_surface_configure(void *data, struct xdg_surface *xdg_surface, uint32_t serial) {
    printf("INFO: wlc_xdg_surface_configure()!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n");
    // client_state *cstate = (client_state *)data;

    // xdg_surface_ack_configure(xdg_surface, serial);
    // wlc_draw_frame(cstate);
    // wl_buffer *buffer = wlc_draw_frame();
}

struct wl_buffer *wlc_draw_frame(client_state *cstate) {
    int retries = 10;

    for (;;) {
        cstate->shmName = wlc_rndShm();
        printf("INFO: Trying %s... ", cstate->shmName);
        cstate->shmFd = shm_open(cstate->shmName, O_RDWR | O_CREAT | O_EXCL, 0600);
        if (cstate->shmFd < 0 && retries > 0) {
            free(cstate->shmName);
            retries--;
            printf("Failed\n");
            continue;
        }
        break;
    }
    if (cstate->shmFd < 0) {
        return NULL;
    }
    printf("Success\n");

    shm_unlink(cstate->shmName);
    free(cstate->shmName);
    // cstate->shm_pool = wl_shm_create_pool(cstate->shm, fd, size);
    return NULL;
}

void wlc_xdg_toplevel_configure(void *data, struct xdg_toplevel *xdg_toplevel, int32_t width, int32_t height, struct wl_array *states) {
    client_state *cstate = (client_state *)data;
    if (width == 0 || height == 0) {
        return;
    }
    cstate->width = width;
    cstate->height = height;
}

void wlc_xdg_toplevel_close(void *data, struct xdg_toplevel *xdg_toplevel) {
    client_state *cstate = (client_state *)data;
    cstate->closed = true;
}

int main(void) {
    int ret = 0;

    client_state cstate = {
        .display = NULL,
        .registry = NULL,
        .compositor = NULL,
        .surface = NULL,
        .shm = NULL,
        .shm_pool = NULL,

        .registry_listener = {
            .global = wlc_registry_global,
            .global_remove = wlc_registry_global_remove
        },

        .xdg_wm_base = NULL,
        .xdg_surface = NULL,
        .xdg_toplevel = NULL,

        .xdg_surface_listener = {
            .configure = wlc_xdg_surface_configure
        },
        .xdg_toplevel_listener = {
            .configure = wlc_xdg_toplevel_configure,
            .close = wlc_xdg_toplevel_close
        },

        .shmName = NULL,
        .shmFd = 0,
        .width = 0,
        .height = 0,
        .closed = false
    };

    cstate.display = wl_display_connect(NULL);
    if (cstate.display == NULL) {
        fprintf(stderr, "ERROR: wl_display_connect()\n");
        return EXIT_FAILURE;
    }

    cstate.registry = wl_display_get_registry(cstate.display);
    if (cstate.registry == NULL) {
        fprintf(stderr, "ERROR: wl_display_get_registry()\n");
        return EXIT_FAILURE;
    }

    // Add listener for events emitted by the registry object
    wl_registry_add_listener(cstate.registry, &cstate.registry_listener, (void *)&cstate);

    ret = wl_display_roundtrip(cstate.display);
    printf("INFO: Dispatched %d events\n", ret);

    if (cstate.compositor == NULL) {
        fprintf(stderr, "error: compositor == NULL\n");
        return EXIT_FAILURE;
    }

    cstate.surface = wl_compositor_create_surface(cstate.compositor);
    if (cstate.surface == NULL) {
        fprintf(stderr, "ERROR: surface == NULL\n");
        return EXIT_FAILURE;
    }

    if (cstate.xdg_wm_base == NULL) {
        fprintf(stderr, "ERROR: xdg_wm_base == NULL\n");
        return EXIT_FAILURE;
    }

    cstate.xdg_surface = xdg_wm_base_get_xdg_surface(cstate.xdg_wm_base, cstate.surface);
    if (cstate.xdg_surface == NULL) {
        fprintf(stderr, "ERROR: xdg_surface == NULL\n");
        return EXIT_FAILURE;
    }

    if (cstate.shm == NULL) {
        fprintf(stderr, "ERROR: shm == NULL\n");
        return EXIT_FAILURE;

    }

    xdg_surface_add_listener(cstate.xdg_surface, &cstate.xdg_surface_listener, &cstate);

    cstate.xdg_toplevel = xdg_surface_get_toplevel(cstate.xdg_surface);
    // xdg_toplevel_add_listener(cstate.xdg_toplevel, )
    // sleep(3);
    wl_display_disconnect(cstate.display);
    return 0;
}
