import apiAuth from '@/api/auth.js'
import apiRole from '@/api/role.js'
import apiUser from '@/api/user.js'
import apiPermission from '@/api/permission.js'

import { basePermissions } from '@/settings'

export async function getUserInfo() {
  const res = await apiAuth.getInfo()
  console.log(res)
  const { id, username, profile, roles, currentRole } = res.data || {}
  return {
    id,
    username,
    avatar: profile?.avatar,
    nickName: profile?.nickName,
    gender: profile?.gender,
    address: profile?.address,
    email: profile?.email,
    roles,
    currentRole,
  }
}

export async function getPermissions() {
  let asyncPermissions = []
  try {
    const res = await apiPermission.getAllPermissionTree()
    asyncPermissions = res?.data || []
  }
  catch (error) {
    console.error(error)
  }
  return basePermissions.concat(asyncPermissions)
}
