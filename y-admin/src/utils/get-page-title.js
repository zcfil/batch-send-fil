import defaultSettings from '@/settings'

const title = defaultSettings.title || '分币系统'

export default function getPageTitle(pageTitle) {
  if (pageTitle) {
    return `${pageTitle} - ${title}`
  }
  return `${title}`
}
