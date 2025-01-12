package com.kuroko.filters;

import java.io.IOException;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.core.annotation.Order;
import org.springframework.stereotype.Component;
import jakarta.servlet.*;
import jakarta.servlet.http.*;

@Component
@Order(1)
public class RequestResponseLoggingFilter implements Filter {
        private final static Logger LOG = LoggerFactory.getLogger(RequestResponseLoggingFilter.class);

        @Override
        public void init(final FilterConfig filterConfig) throws ServletException {
                LOG.info("Initializing filter :{}", this);
        }

        @Override
        public void doFilter(
                        ServletRequest request,
                        ServletResponse response,
                        FilterChain chain) throws ServletException, IOException {

                HttpServletRequest req = (HttpServletRequest) request;
                HttpServletResponse res = (HttpServletResponse) response;
                LOG.info(
                                "Logging Request  {} : {}", req.getMethod(),
                                req.getRequestURI());
                chain.doFilter(request, response);
                LOG.info(
                                "Logging Response :{}",
                                res.getContentType());
        }

        // other methods
}
