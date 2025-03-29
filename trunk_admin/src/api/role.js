import { request } from '@/utils'


export default {
    create: data => request.post('/role', data),
    read: (params = {}) => request.get('/role/page', { params }),
    update: data => request.patch(`/role/${data.id}`, data),
    delete: id => request.delete(`/role/${id}`),
  
    // 获取角色权限
    getAllRoles: () => request.get('/role?enable=1'),
    addRoleUsers: (roleId, data) => request.patch(`/role/users/add/${roleId}`, data),
    removeRoleUsers: (roleId, data) => request.patch(`/role/users/remove/${roleId}`, data),
  }
