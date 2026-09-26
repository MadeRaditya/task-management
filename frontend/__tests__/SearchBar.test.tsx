import { fireEvent, render } from '@testing-library/react-native';
import { SearchBar } from '../src/components/SearchBar';

describe('SearchBar', () => {
  it('memanggil onChangeText saat user mengetik', () => {
    const onChangeText = jest.fn();
    const { getByTestId } = render(
      <SearchBar value="" onChangeText={onChangeText} />
    );

    fireEvent.changeText(getByTestId('search-input'), 'meeting');
    expect(onChangeText).toHaveBeenCalledWith('meeting');
  });

  it('menampilkan value yang diberikan', () => {
    const { getByTestId } = render(
      <SearchBar value="hello" onChangeText={() => {}} />
    );
    expect(getByTestId('search-input').props.value).toBe('hello');
  });
});