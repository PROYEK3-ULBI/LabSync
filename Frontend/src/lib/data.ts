// Data contoh fiktif untuk demonstrasi antarmuka — bukan data institusi resmi.

export const currentStudent = {
  name: 'Ahmad Fauzan',
  nim: 'TI240001',
  email: 'ahmad.fauzan@student.univ.ac.id',
  prodi: 'Teknik Informatika',
  angkatan: '2024',
  status: 'Aktif',
  faceEnrolled: true,
}

export const activeSession = {
  code: 'PWB-07',
  name: 'Praktikum Pemrograman Web',
  course: 'Pemrograman Web',
  date: '18 September 2026',
  start: '13.00',
  end: '15.30',
  room: 'Lab Komputer 1',
  status: 'Aktif',
  attendance: 'Belum Hadir',
  present: 22,
  total: 30,
}

export type Session = {
  code: string
  name: string
  course: string
  date: string
  start: string
  end: string
  room: string
  status: string
  attendance?: string
}

export const sessions: Session[] = [
  { code: 'PWB-07', name: 'Praktikum Pemrograman Web', course: 'Pemrograman Web', date: '18 Sep 2026', start: '13.00', end: '15.30', room: 'Lab Komputer 1', status: 'Aktif', attendance: 'Belum Hadir' },
  { code: 'BDT-05', name: 'Praktikum Basis Data', course: 'Basis Data', date: '19 Sep 2026', start: '08.00', end: '10.30', room: 'Lab Komputer 2', status: 'Akan Datang', attendance: 'Belum Hadir' },
  { code: 'JAR-04', name: 'Praktikum Jaringan Komputer', course: 'Jaringan Komputer', date: '20 Sep 2026', start: '10.00', end: '12.30', room: 'Lab Komputer 1', status: 'Akan Datang', attendance: 'Belum Hadir' },
  { code: 'PWB-06', name: 'Praktikum Pemrograman Web', course: 'Pemrograman Web', date: '11 Sep 2026', start: '13.00', end: '15.30', room: 'Lab Komputer 1', status: 'Ditutup', attendance: 'Hadir' },
  { code: 'BDT-04', name: 'Praktikum Basis Data', course: 'Basis Data', date: '12 Sep 2026', start: '08.00', end: '10.30', room: 'Lab Komputer 2', status: 'Ditutup', attendance: 'Hadir' },
  { code: 'JAR-03', name: 'Praktikum Jaringan Komputer', course: 'Jaringan Komputer', date: '13 Sep 2026', start: '10.00', end: '12.30', room: 'Lab Komputer 1', status: 'Dibatalkan', attendance: 'Tidak Hadir' },
]

export type AttendanceRow = {
  date: string
  session: string
  room: string
  checkIn: string
  checkOut: string
  status: string
}

export const attendanceHistory: AttendanceRow[] = [
  { date: '11 Sep 2026', session: 'Praktikum Pemrograman Web', room: 'Lab Komputer 1', checkIn: '13.02', checkOut: '15.28', status: 'Selesai' },
  { date: '12 Sep 2026', session: 'Praktikum Basis Data', room: 'Lab Komputer 2', checkIn: '08.05', checkOut: '10.30', status: 'Selesai' },
  { date: '04 Sep 2026', session: 'Praktikum Pemrograman Web', room: 'Lab Komputer 1', checkIn: '13.00', checkOut: '15.30', status: 'Selesai' },
  { date: '05 Sep 2026', session: 'Praktikum Basis Data', room: 'Lab Komputer 2', checkIn: '08.10', checkOut: '—', status: 'Pending' },
  { date: '30 Agu 2026', session: 'Praktikum Jaringan Komputer', room: 'Lab Komputer 1', checkIn: '—', checkOut: '—', status: 'Gagal' },
]

export const myAsset = {
  code: 'LAB-LT-007',
  name: 'Dell Latitude 5490 Core i7',
  brand: 'Dell',
  model: 'Latitude 5490',
  session: 'Praktikum Pemrograman Web',
  allocatedAt: '18 Sep 2026, 13.04',
  condition: 'Baik',
  status: 'Allocated',
}

export const assetHistory = [
  { date: '11 Sep 2026', code: 'LAB-LT-007', session: 'Pemrograman Web', out: '13.05', in: '15.29', condition: 'Baik', status: 'Selesai' },
  { date: '04 Sep 2026', code: 'LAB-LT-012', session: 'Pemrograman Web', out: '13.02', in: '15.31', condition: 'Baik', status: 'Selesai' },
  { date: '28 Agu 2026', code: 'LAB-LT-007', session: 'Basis Data', out: '08.06', in: '10.33', condition: 'Perlu Perbaikan', status: 'Selesai' },
]

export type User = {
  name: string
  nim: string
  email: string
  role: string
  status: string
  face: string
  created: string
}
export const users: User[] = [
  { name: 'Ahmad Fauzan', nim: 'TI240001', email: 'ahmad.fauzan@student.univ.ac.id', role: 'Mahasiswa', status: 'Aktif', face: 'Terdaftar', created: '02 Sep 2026' },
  { name: 'Siti Rahma', nim: 'TI240002', email: 'siti.rahma@student.univ.ac.id', role: 'Mahasiswa', status: 'Aktif', face: 'Terdaftar', created: '02 Sep 2026' },
  { name: 'Muhammad Rizky', nim: 'TI240003', email: 'm.rizky@student.univ.ac.id', role: 'Mahasiswa', status: 'Aktif', face: 'Belum Terdaftar', created: '03 Sep 2026' },
  { name: 'Nabila Putri', nim: 'TI240004', email: 'nabila.putri@student.univ.ac.id', role: 'Mahasiswa', status: 'Nonaktif', face: 'Terdaftar', created: '03 Sep 2026' },
  { name: 'Dewi Anggraini', nim: 'LAB0001', email: 'dewi.laboran@univ.ac.id', role: 'Laboran', status: 'Aktif', face: '—', created: '20 Agu 2026' },
  { name: 'Budi Santoso', nim: 'TU0002', email: 'budi.tu@univ.ac.id', role: 'Tata Usaha', status: 'Aktif', face: '—', created: '20 Agu 2026' },
]

export type Asset = {
  code: string
  name: string
  category: string
  brand: string
  model: string
  serial: string
  location: string
  status: string
}
export const assets: Asset[] = [
  { code: 'LAB-LT-007', name: 'Dell Latitude 5490 Core i7', category: 'Laptop', brand: 'Dell', model: 'Latitude 5490', serial: 'DL5490-0091', location: 'Lab Komputer 1', status: 'Allocated' },
  { code: 'LAB-LT-008', name: 'Dell Latitude 5490 Core i5', category: 'Laptop', brand: 'Dell', model: 'Latitude 5490', serial: 'DL5490-0092', location: 'Lab Komputer 1', status: 'Available' },
  { code: 'LAB-PC-021', name: 'HP ProDesk 400 G7', category: 'Desktop', brand: 'HP', model: 'ProDesk 400 G7', serial: 'HP400-0210', location: 'Lab Komputer 2', status: 'Available' },
  { code: 'LAB-PC-022', name: 'HP ProDesk 400 G7', category: 'Desktop', brand: 'HP', model: 'ProDesk 400 G7', serial: 'HP400-0211', location: 'Lab Komputer 2', status: 'Maintenance' },
  { code: 'LAB-SW-003', name: 'Cisco Catalyst 2960', category: 'Perangkat Jaringan', brand: 'Cisco', model: 'Catalyst 2960', serial: 'CS2960-0033', location: 'Lab Komputer 1', status: 'Available' },
  { code: 'LAB-LT-013', name: 'Lenovo ThinkPad E14', category: 'Laptop', brand: 'Lenovo', model: 'ThinkPad E14', serial: 'LN-E14-0130', location: 'Gudang', status: 'Retired' },
]

export type Allocation = {
  id: string
  student: string
  nim: string
  session: string
  asset: string
  assetName: string
  time: string
  condition: string
  status: string
}
export const allocations: Allocation[] = [
  { id: 'ALC-1042', student: 'Ahmad Fauzan', nim: 'TI240001', session: 'Pemrograman Web', asset: 'LAB-LT-007', assetName: 'Dell Latitude 5490', time: '18 Sep 2026, 13.04', condition: 'Baik', status: 'Aktif' },
  { id: 'ALC-1041', student: 'Siti Rahma', nim: 'TI240002', session: 'Pemrograman Web', asset: 'LAB-PC-021', assetName: 'HP ProDesk 400 G7', time: '18 Sep 2026, 13.06', condition: 'Baik', status: 'Aktif' },
  { id: 'ALC-1038', student: 'Muhammad Rizky', nim: 'TI240003', session: 'Basis Data', asset: 'LAB-LT-008', assetName: 'Dell Latitude 5490', time: '12 Sep 2026, 08.10', condition: 'Baik', status: 'Selesai' },
]

export type Maintenance = {
  id: string
  code: string
  name: string
  issue: string
  priority: string
  status: string
  reported: string
  completed: string
}
export const maintenance: Maintenance[] = [
  { id: 'MT-208', code: 'LAB-PC-022', name: 'HP ProDesk 400 G7', issue: 'Tidak dapat booting', priority: 'Tinggi', status: 'Sedang Dikerjakan', reported: '17 Sep 2026', completed: '—' },
  { id: 'MT-205', code: 'LAB-LT-011', name: 'Dell Latitude 5490', issue: 'Keyboard rusak', priority: 'Sedang', status: 'Terbuka', reported: '16 Sep 2026', completed: '—' },
  { id: 'MT-201', code: 'LAB-SW-002', name: 'Cisco Catalyst 2960', issue: 'Port 12 mati', priority: 'Rendah', status: 'Selesai', reported: '10 Sep 2026', completed: '14 Sep 2026' },
]

export const attendanceMonitoring = [
  { student: 'Ahmad Fauzan', nim: 'TI240001', session: 'Pemrograman Web', checkIn: '13.02', checkOut: '—', verify: 'Terverifikasi', status: 'Hadir' },
  { student: 'Siti Rahma', nim: 'TI240002', session: 'Pemrograman Web', checkIn: '13.06', checkOut: '—', verify: 'Terverifikasi', status: 'Hadir' },
  { student: 'Muhammad Rizky', nim: 'TI240003', session: 'Pemrograman Web', checkIn: '—', checkOut: '—', verify: 'Gagal', status: 'Tidak Hadir' },
  { student: 'Nabila Putri', nim: 'TI240004', session: 'Pemrograman Web', checkIn: '—', checkOut: '—', verify: '—', status: 'Tidak Hadir' },
]

export const auditLog = [
  { time: '18 Sep 2026 13.04', user: 'Ahmad Fauzan', action: 'Face Verification', entity: 'Presensi PWB-07', status: 'Terverifikasi', ip: '10.20.3.41' },
  { time: '18 Sep 2026 13.04', user: 'Ahmad Fauzan', action: 'Attendance Created', entity: 'PWB-07', status: 'Selesai', ip: '10.20.3.41' },
  { time: '18 Sep 2026 13.02', user: 'Dewi Anggraini', action: 'Allocation Created', entity: 'ALC-1042', status: 'Selesai', ip: '10.20.3.10' },
  { time: '18 Sep 2026 12.50', user: 'Dewi Anggraini', action: 'Asset Created', entity: 'LAB-SW-003', status: 'Selesai', ip: '10.20.3.10' },
  { time: '17 Sep 2026 16.20', user: 'Budi Santoso', action: 'Maintenance Created', entity: 'MT-208', status: 'Selesai', ip: '10.20.3.12' },
  { time: '17 Sep 2026 08.31', user: 'Muhammad Rizky', action: 'Login', entity: 'Sesi Login', status: 'Gagal', ip: '10.20.4.77' },
]

export const recentActivity = auditLog.slice(0, 5)

// Weekly attendance for a small bar chart
export const weeklyAttendance = [
  { day: 'Sen', hadir: 26, total: 30 },
  { day: 'Sel', hadir: 28, total: 30 },
  { day: 'Rab', hadir: 24, total: 28 },
  { day: 'Kam', hadir: 29, total: 30 },
  { day: 'Jum', hadir: 22, total: 30 },
]
