import {defineStore} from "pinia";
import {Credentials} from "components/models";
import AuthService from "src/service/AuthService";

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: '',
    username: '',
  }),
  getters: {
    isAuthenticated(): boolean {
      return this.token !== '';
    },
  },
  actions: {
    logout() {
      this.token = '';
      this.username = '';
    },
    async login(credentials: Credentials): Promise<void> {
      return new Promise((resolve, reject) => {
        AuthService.login(credentials).then((res) => {
          resolve();
        }).catch((err) => reject(err));
      });
    },
    async register(credentials: Credentials): Promise<void> {
      return new Promise((resolve, reject) => {
        AuthService.register(credentials).then(() => {
          resolve();
        }).catch((err) => reject(err));
      });
    }
  }
});
