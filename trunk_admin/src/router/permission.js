import router from '.'
import NProgress from 'nprogress' // Progress 进度条
import 'nprogress/nprogress.css'// Progress 进度条样式
 
NProgress.configure({
  easing: 'ease', // 动画方式    
  speed: 500, // 递增进度条的速度    
  showSpinner: false, // 是否显示加载ico    
  trickleSpeed: 200, // 自动递增间隔    
  minimum: 0.3 // 初始化时的最小百分比
})



import StoreAuth from '@/store/auth.js'
import apiAuth from '@/api/auth'





let whitelist = ["/about/index", "/login", "/register","/404"]


router.beforeEach(async (to, from, next) => {
  const storeAuth = StoreAuth()

  NProgress.start(); // 开启Progress

  if (whitelist.indexOf(to.path) !== -1) {
    next()
    return
  }
  
  next()
  
  /*
  try {
    let ret = await apiAuth.info();
    if (ret.code == 200) {
      storeAuth.setAuthInfo(ret.data)
      next()
      return;
    }
    else if (ret.code == 401) {
      next('/login')
      return;
    }

  } catch (error) {
    console.log("router.beforeEach： error:", error)
    next('/login')
  }
  console.log("router beforeEach stop  ");
  next(false);
  */

});

router.afterEach(() => {
  NProgress.done() // 结束Progress
});


