import { defineStore } from 'pinia' //引入
import apiAuth from '@/api/auth'
import apiUser from '@/api/user'

const PiniaStore = defineStore('storeAuth', {
  state: () => {
    return {
        //预置字段
        userInfo:{},      //用户信息
        authInfo:{},      //认证信息
        AdminOption:false,//管理员操作
    }
  },
  persist:true,
  getters: {

    Logined(){
      return this.authInfo != {};
    },

    Token(store){
      return store.authInfo.token;
    },
    IsAdmin(store){
      return store.userInfo.is_admin;
    },
  },
  actions: {

        //设置认证信息--包含了token
        setAuthInfo(info) {
            this.authInfo = JSON.parse(JSON.stringify(info));
        },

        // 设置用户信息
        setUserInfo(info) {
            this.userInfo = JSON.parse(JSON.stringify(info));
        },
        

        //切换管理员操作
        switchAdminOption(value) {
            if(this.IsAdmin) {
                this.AdminOption = value;
            }else{
                console.log("你不是管理员,无法切换到管理员操作");
            }
        },
        // 新增方法：定时获取用户信息
        startFetchingAuthInfo(interval) {
            setInterval(() => {
                if(this.Logined) {
                    this.fetchAuthInfo();
                }
            }, interval);
        },
      // 假设你有一个方法来获取用户信息
        async fetchAuthInfo() {
            try {
                const ret = await apiAuth.info();
                if (ret.code == 200) {
                    //console.log('认证信息:', ret.data);
                    this.setAuthInfo(ret.data);
                }
            } catch (error) {
                console.error('Failed to fetch user info:', error);
            }
        },
  },
  //persist:true,
})



export default PiniaStore //导出
