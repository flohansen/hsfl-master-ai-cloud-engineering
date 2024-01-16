<template>
  <q-layout view="hHh Lpr lff">
    <q-header elevated class="bg-primary text-white">
      <q-toolbar>
        <q-btn dense flat round icon="menu" @click="toggleLeftDrawer" />

        <q-toolbar-title>
          <span>BoardHub</span>
        </q-toolbar-title>

        <q-btn color="black" label="USERNAME">
          <q-menu>
            <q-list style="min-width: 100px">
              <q-separator />
              <q-item class="text-red" clickable v-close-popup>
                <q-item-section @click="logout">Logout</q-item-section>
              </q-item>
            </q-list>
          </q-menu>
        </q-btn>
      </q-toolbar>
    </q-header>

    <q-drawer
      v-model="leftDrawerOpen"
      show-if-above
      :width="200"
      :breakpoint="500"
      bordered
      :class="$q.dark.isActive ? 'bg-grey-9' : 'bg-grey-3'"
    >
      <q-scroll-area class="fit">
        <q-list>

          <template v-for="(menuItem, index) in menuList" :key="index">
            <q-item clickable :active="menuItem.label === 'Outbox'" v-ripple @click="goTo(menuItem.route)">
              <q-item-section avatar>
                <q-icon :name="menuItem.icon" />
              </q-item-section>
              <q-item-section>
                {{ menuItem.label }}
              </q-item-section>
            </q-item>
            <q-separator :key="'sep' + index"  v-if="menuItem.separator" />
          </template>

        </q-list>
      </q-scroll-area>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import {ref} from 'vue';
import {useRouter} from 'vue-router';
import {useAuthStore} from 'stores/auth-store';

const menuList = [
  {
    route:'home',
    icon: 'dashboard',
    label:'Home',
    separator: false
  }
]

const router = useRouter()
const authStore = useAuthStore()
const leftDrawerOpen = ref(false)

function goTo(route: string) {
  leftDrawerOpen.value = false
  router.push({ name: route })
}

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value
}

function logout() {
  authStore.logout();
}
</script>
