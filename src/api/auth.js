// import Config from '../config.js';

export const login = (username, password) => {
  console.log('Received login request ');
  console.log(username);
  console.log(password);
  return(
    {
      success: true,
      token: 'aaaabbbbccccddddeeee',
    }
  );
}
