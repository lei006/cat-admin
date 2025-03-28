import { request } from '@/utils'




export default {
    // 获取用户信息
    login: data => request.post('/auth/login', data, { needToken: false }),

    // 退出登录
    logout: () => request.post('/auth/logout', {}, { needTip: false }),

    // 刷新token
    refreshToken: () => request.get('/auth/refresh/token'),

    // 取得当前用户信息
    getUser: () => request.get('/user/detail'),

    // 切换角色
    toggleRole: data => request.post('/auth/role/toggle', data),

    // 获取角色权限
    getRolePermissions: () => request.get('/role/permissions/tree'),

    // 切换当前角色
    switchCurrentRole: role => request.post(`/auth/current-role/switch/${role}`),

  }

  
