import { fireEvent, render } from '@testing-library/react-native';
import { TaskList } from '../src/components/TaskList';
import { Task } from '../src/types/task';

const sample: Task[] = [
  {
    id: 1,
    title: 'Belajar Go',
    description: 'Bikin API',
    status: 'todo',
    assignee: 'Radit',
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
  },
  {
    id: 2,
    title: 'Setup Redis',
    description: 'Cache GET',
    status: 'done',
    assignee: 'Sinta',
    created_at: '2026-01-02T00:00:00Z',
    updated_at: '2026-01-02T00:00:00Z',
  },
];

describe('TaskList', () => {
  it('menampilkan semua task', () => {
    const { getByText } = render(
      <TaskList tasks={sample} onEdit={() => {}} onDelete={() => {}} />
    );
    expect(getByText('Belajar Go')).toBeTruthy();
    expect(getByText('Setup Redis')).toBeTruthy();
  });

  it('memanggil onEdit saat tombol edit ditekan', () => {
    const onEdit = jest.fn();
    const { getByTestId } = render(
      <TaskList tasks={sample} onEdit={onEdit} onDelete={() => {}} />
    );
    fireEvent.press(getByTestId('task-edit-1'));
    expect(onEdit).toHaveBeenCalledWith(sample[0]);
  });

  it('memanggil onDelete saat tombol hapus ditekan', () => {
    const onDelete = jest.fn();
    const { getByTestId } = render(
      <TaskList tasks={sample} onEdit={() => {}} onDelete={onDelete} />
    );
    fireEvent.press(getByTestId('task-delete-2'));
    expect(onDelete).toHaveBeenCalledWith(sample[1]);
  });

  it('menampilkan pesan kosong saat tasks kosong', () => {
    const { getByText } = render(
      <TaskList tasks={[]} onEdit={() => {}} onDelete={() => {}} />
    );
    expect(getByText('Tidak ada task.')).toBeTruthy();
  });
});