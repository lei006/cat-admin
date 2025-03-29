import { request } from '@/utils'

export default {
    create: data => request.post('/user', data),
    read: (params = {}) => request.get('/user', { params }),
    update: data => request.patch(`/user/${data.id}`, data),
    delete: id => request.delete(`/user/${id}`),
    resetPwd: (id, data) => request.patch(`/user/password/reset/${id}`, data),
  

    // 取得当前用户信息
    getAllUsers: (params = {}) => request.get('/user', { params }),

    updateProfile: data => request.patch(`/user/profile/${data.id}`, data),

  }

  