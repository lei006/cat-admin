import { request } from '@/utils'
import axios from 'axios'

export default {
  getMenuTree: () => request.get('/permission/menu/tree'),
  getButtons: ({ parentId }) => request.get(`/permission/button/${parentId}`),
  getComponents: () => axios.get(`${import.meta.env.VITE_PUBLIC_PATH}components.json`),
  addPermission: data => request.post('/permission', data),
  savePermission: (id, data) => request.patch(`/permission/${id}`, data),
  deletePermission: id => request.delete(`permission/${id}`),
  getAllPermissionTree: () => request.get('/permission/tree'),
  
  // 验证菜单路径
  validateMenuPath: path => request.get(`/permission/menu/validate?path=${path}`),  
}

