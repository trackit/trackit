import React, { Component } from 'react';
import PropTypes from 'prop-types';

class ResultComponent extends Component {
    getFACode(status) {
        switch (status) {
            case 'green':
                return 'fa-check-circle';        
            case 'orange':
                return 'fa-exclamation-circle';        
            case 'red':
                return 'fa-exclamation-triangle';        
            default:
                return '';        
        }
    }

    getPercentage(passed, checked) {
        if (checked === 0) {
            return 100;
        }
        return ((passed * 100) / checked);
    }

    render() {
        const { result } = this.props;

        let details;
        if (result.details && result.details.length) {
            details = (
                <div>
                    <br />
                    <div className="row">
                        <div className="col-md-7 col-md-offset-1">
                            <pre className="details-pre">
                                <ul>
                                    {result.details.map(item => <li key={item}>{item}</li>)}
                                </ul>
                            </pre>
                        </div>
                    </div>
                </div>
            );    
        }

        const resultBlock = (
            <div>
                <div className="row">
                    <div className="col-md-1 icon-col">
                        <i className={`fa fa-2x ${this.getFACode(result.status)} ${result.status}-color`}></i>
                    </div>
                    <div className="col-md-7">
                        <strong>
                            {`${result.category} - ${result.pluginName}`}
                        </strong>
                        <br />
                        <p className="no-padding no-margin">{result.result}</p>
                    </div>
                    <div className="col-md-3">
                        <div className="progress">
                            <div
                                className={`progress-bar ${result.status}-bg`}
                                role="progressbar"
                                aria-valuenow={this.getPercentage(result.passed, result.checked)}
                                aria-valuemin="0"
                                aria-valuemax="100"
                                style={{width: `${this.getPercentage(result.passed, result.checked)}%`}}
                            >
                            </div>
                        </div>

                        <p className="no-padding no-margin">
                            {`${result.label}: ${result.passed} / ${result.checked}`}
                        </p>
                    </div>
                </div>
                {details}
            </div>
        );

        const errorBlock = (
            <div className="row">
                <div className="col-md-8 col-md-offset-1">
                    <div className="alert alert-warning">
                        <strong>{result.pluginName} - Error :</strong>
                        &nbsp;
                        {result.error}
                    </div>
                </div>
            </div>
        );

        return (
            <div className="conformity-result">
                {result.error.length ? errorBlock : resultBlock}
            </div>
        );
    }
}

ResultComponent.propTypes = {
    result: PropTypes.object.isRequired,
}

export default ResultComponent;