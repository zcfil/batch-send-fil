import Vue from 'vue'


Vue.filter('NumberFormat', function (s) {
  if (!s) {
    return '0'
  }
  s = s + ''
  // eslint-disable-next-line no-useless-escape
  if (/[^0-9\.]/.test(s)) return ''
  s = s.replace(/^(\d*)$/, '$1.')
  s = (s + '00').replace(/(\d*\.\d\d)\d*/, '$1')
  s = s.replace('.', ',')
  const re = /(\d)(\d{3},)/
  while (re.test(s)) {
    s = s.replace(re, '$1,$2')
  }
  s = s.replace(/,(\d\d)$/, '.$1')
  return s.replace(/^\./, '0.')
})

