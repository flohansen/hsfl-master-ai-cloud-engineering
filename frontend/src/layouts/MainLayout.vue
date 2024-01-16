<template>
  <q-layout view="hHh Lpr lff">
    <q-header elevated class="bg-primary text-white">
      <q-toolbar>
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

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import {useRouter} from 'vue-router';
import {useAuthStore} from 'stores/auth-store';

const router = useRouter()
const authStore = useAuthStore()

function logout() {
  authStore.logout();
  router.push({name: 'login'})
}
</script>
