<template>
  <q-page class="row items-center justify-evenly">
    <q-infinite-scroll @load="onLoad" :offset="500" :disable="true">
      <div v-for="(entry, index) in bulletinStore.get" :key="index" class="caption">
        <q-card bordered class="q-ma-lg" @click="() => {
          displayBulletinBoardEntry = entry
        }"
        style="min-width: 400px">
          <q-card-section>
            <div class="text-h6">{{ entry.title }}</div>
            <div class="text-subtitle2">at {{ entry.createdAt }}</div>
          </q-card-section>

          <q-separator dark inset />

          <q-card-section>
            {{ entry.description }}
          </q-card-section>
        </q-card>
      </div>
      <template v-slot:loading>
        <div class="row justify-center q-my-md">
          <q-spinner-dots color="primary" size="40px" />
        </div>
      </template>
    </q-infinite-scroll>
    <q-page-sticky position="bottom-right" :offset="[18, 18]">
      <q-btn fab icon="add" color="accent" @click="displayCreateDialog = true" />
    </q-page-sticky>

    <!-- Create Dialog -->
    <CreateBulletinDialog :visible="displayCreateDialog" @close="displayCreateDialog = false" @create="bulletinStore.create" />

    <!-- Display Dialog -->
    <DisplayBulletinDialog :item="displayBulletinBoardEntry" @close="displayBulletinBoardEntry = null" />

  </q-page>
</template>

<script setup lang="ts">
import CreateBulletinDialog from "components/bulletin-board/CreateBulletinDialog.vue";


import {onMounted, ref} from "vue";
import {BulletinBoardEntry} from "components/models";
import DisplayBulletinDialog from "components/bulletin-board/DisplayBulletinDialog.vue";
import {useBulletinStore} from "stores/bulletin-store";

const bulletinStore = useBulletinStore()

const displayCreateDialog = ref(false)
const displayBulletinBoardEntry= ref<BulletinBoardEntry | null>(null)

const onLoad = (index: any, done: () => void) => {
  setTimeout(() => {
    done()
  }, 500)
}

onMounted(() => {
  //bulletinStore.fetch()
})

</script>
