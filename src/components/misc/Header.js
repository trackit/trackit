import React, {Component} from 'react';

import logo from '../../assets/logo.png';

class Header extends Component {

  render() {
    const styles = {
      container: {
        marginBottom: '30px',
        padding: '10px',
        backgroundColor: '#d6413b',
        color: 'white',
      },
      logo: {
        maxWidth: '200px',
        maxHeight: '50px',
        marginLeft: '30px',
      },
      title: {
        color: 'white',
        fontSize: '28px',
        margin: '10px auto',
      }

    }

    return(
      <div className="text-center" style={styles.container}>
          <a
            href="https://trackit.io"
            rel="noopener noreferrer"
            target="_blank"
            className="pull-left animated bounceInLeft"
          >
            <img
              src={logo}
              alt="Trackit Markets Cloud Storage Comparator"
              style={styles.logo}
            />
          </a>
          <h1 style={styles.title} className=" animated bounceInRight">
            Cloud Storage Comparator
          </h1>
          <div style={{ clear: 'both' }} />
      </div>
    );
  }

}

export default Header;
