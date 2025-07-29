import moment from 'moment-jalaali';

moment.loadPersian({ dialect: 'persian-modern', usePersianDigits: true });

export function toJalali(dateStr) {
  return moment(dateStr).format('jYYYY/jMM/jDD');
}

export function toGregorian(dateStr) {
  const [jy, jm, jd] = dateStr.split(/[-\/]/).map(Number);
  const g = window.jalaali.toGregorian(jy, jm, jd);
  return `${g.gy}-${String(g.gm).padStart(2, '0')}-${String(g.gd).padStart(2, '0')}`;
}
