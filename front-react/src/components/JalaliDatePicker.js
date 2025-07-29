import React from 'react';
import { DatePicker } from 'react-modern-calendar-datepicker';

export default function JalaliDatePicker({ value, onChange }) {
  // value is a string YYYY-MM-DD in Gregorian
  const gregToJalali = dateStr => {
    if (!dateStr) return null;
    const [gy, gm, gd] = dateStr.split('-').map(Number);
    const j = window.jalaali.toJalaali(gy, gm, gd);
    return { year: j.jy, month: j.jm, day: j.jd };
  };
  const jalaliToGreg = dateObj => {
    if (!dateObj) return '';
    const g = window.jalaali.toGregorian(dateObj.year, dateObj.month, dateObj.day);
    return `${g.gy}-${String(g.gm).padStart(2, '0')}-${String(g.gd).padStart(2, '0')}`;
  };
  const handleChange = d => onChange(jalaliToGreg(d));
  return (
    <DatePicker
      value={gregToJalali(value)}
      onChange={handleChange}
      locale="fa"
      inputPlaceholder="انتخاب تاریخ"
      shouldHighlightWeekends
    />
  );
}
