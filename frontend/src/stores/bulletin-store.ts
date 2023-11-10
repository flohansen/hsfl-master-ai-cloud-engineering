import {defineStore} from "pinia";
import BulletinBoardService from "src/service/BulletinBoardService";
import {BulletinBoardEntry} from "components/models";

export const useBulletinStore = defineStore('bulletin', {
  state: () => ({
    bulletins: [
      {
        id: '1',
        title: 'First bulletin',
        description: 'This is the first bulletin',
        createdAt: '2021-01-01T00:00:00.000Z',
      },
      {
        id: '2',
        title: 'Second bulletin',
        description: 'This is the second bulletin',
        createdAt: '2021-01-02T00:00:00.000Z',
      },
      {
        id: '3',
        title: 'Third bulletin',
        description: 'This is the third bulletin',
        createdAt: '2021-01-03T00:00:00.000Z',
      },
      {
        id: '4',
        title: 'Fourth bulletin',
        description: 'This is the fourth bulletin',
        createdAt: '2021-01-04T00:00:00.000Z',
      },
      {
        id: '5',
        title: 'Fifth bulletin',
        description: 'This is the fifth bulletin',
        createdAt: '2021-01-05T00:00:00.000Z',
      },
      {
        id: '6',
        title: 'Sixth bulletin',
        description: 'This is the sixth bulletin',
        createdAt: '2021-01-06T00:00:00.000Z',
      }
    ] as BulletinBoardEntry[],
  }),
  getters: {
    get(): BulletinBoardEntry[] {
      return this.bulletins;
    },
  },
  actions: {
    fetch() {
        BulletinBoardService.get().then((res) => {
          this.bulletins.length = 0;
          this.bulletins.push(...res);
        });
    },
    create(bulletin: BulletinBoardEntry) {
      BulletinBoardService.create(bulletin).then((res) => {
        this.bulletins.unshift(res);
      });
    },
  }
});
