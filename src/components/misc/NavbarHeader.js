import React, {Component} from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';
import PropTypes from 'prop-types';
import Actions from '../../actions';
import SelectedIndicator from '../aws/accounts/SelectedIndicatorComponent';
import '../../styles/Navigation.css';

import logo from '../../assets/logo-white-coloured.png';

export class NavbarHeader extends Component {

  constructor() {
    super();
    this.state = {
      userMenuExpanded: false,
    };
  }

  toggleUserMenu() {
    this.setState({ userMenuExpanded: !this.state.userMenuExpanded });
  }

  render() {

    return(
      <div>

        <nav className="navbar navbar-default navbar-static-top">
            <div className="navbar-header">

              <div className="top-left-part">
                <Link to='/app' className="logo">
                  <b className="animated fadeInDown">
                    <img src={logo} alt="Trackit logo"/>
                  </b>
                </Link>
              </div>

              <div className="top-right-part pull-right">
                 <span style={{ display: 'inline-block', marginTop: '19px' }}><SelectedIndicator longVersion={true} icon={true} /></span>
                 <div className={this.state.userMenuExpanded ? 'open dropdown-trigger' : 'dropdown-trigger'}>
                   <button className="navbar-user-dropdown-toggle" onClick={this.toggleUserMenu.bind(this)}>
                     <span className="fa-stack red-color">
                       <i className="fa fa-circle fa-stack-2x"></i>
                       <i className="fa fa-user fa-stack-1x fa-inverse"></i>
                     </span>
                     <i className="fa fa-caret-down" />
                   </button>
                   <ul className="dropdown-menu">
                     <li>
                       <h5 className="dropdown-title"><strong>{this.props.mail}</strong></h5>
                     </li>
                     <hr className="m-b-0"/>
                     <li>
                       <Link to='/app/setup'>
                         <i className="fa fa-cog"/>
                         &nbsp;Setup
                       </Link>
                     </li>
                     <li>
                       <a href="" onClick={this.props.signOut}>
                         <i className="fa fa-sign-out"/>
                         &nbsp;Sign out
                       </a>
                     </li>
                   </ul>

                 </div>
              </div>

            </div>
        </nav>

      </div>
    );
  }

}

NavbarHeader.propTypes = {
  mail: PropTypes.string,
  signOut: PropTypes.func.isRequired
};

/* istanbul ignore next */
const mapStateToProps = (state) => ({
  mail: state.auth.mail
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  signOut: () => {
    dispatch(Actions.Auth.logout())
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(NavbarHeader);
