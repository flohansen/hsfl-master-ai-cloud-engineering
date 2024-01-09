<script setup lang="ts">
import {ref} from "vue";
import {useRouter} from "vue-router";

const router = useRouter()

  const email = ref('')
  const password = ref('')
  const passwordRepeated = ref('')

  const register = () => {
    console.log('register')
  }

  const goToLogin = () => {
    router.push({name: 'login'})
  }
</script>

<template>
      <q-page class="flex flex-center bg-grey-2">
        <q-card class="q-pa-md shadow-2 my_card" bordered>
          <q-card-section class="text-center">
            <div class="text-grey-9 text-h5 text-weight-bold">Register</div>
            <div class="text-grey-8">Register below to access the application</div>
          </q-card-section>
          <q-card-section>
            <q-input dense outlined v-model="email" type="email" label="Email Address"></q-input>
            <q-input dense outlined class="q-mt-md" v-model="password" type="password" label="Password"></q-input>
            <q-input dense outlined class="q-mt-md" v-model="passwordRepeated" type="password" label="Repeat Password"></q-input>
          </q-card-section>
          <q-card-section>
            <q-btn style="
  border-radius: 8px;" color="dark" rounded size="md" label="Register" no-caps class="full-width"></q-btn>
          </q-card-section>
          <q-card-section class="text-center q-pt-none">
            <div class="text-grey-8">Already have an account?
              <a @click="goToLogin" class="text-dark text-weight-bold" style="text-decoration: none; cursor: pointer">Sign
                in.</a></div>
          </q-card-section>

        </q-card>
      </q-page>
</template>

<style scoped>
.my_card {
  width: 25rem;
  border-radius: 8px;
  box-shadow: 0 20px 25px -5px rgb(0 0 0 / 0.1), 0 8px 10px -6px rgb(0 0 0 / 0.1);
}
</style>
